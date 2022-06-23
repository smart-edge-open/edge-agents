// INTEL CONFIDENTIAL
//
// Copyright 2021-2021 Intel Corporation.
//
// This software and the related documents are Intel copyrighted materials, and your use of
// them is governed by the express license under which they were provided to you ("License").
// Unless the License provides otherwise, you may not use, modify, copy, publish, distribute,
// disclose or transmit this software or the related documents without Intel's prior written permission.
//
// This software and the related documents are provided as is, with no express or implied warranties,
// other than those that are expressly stated in the License.

package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"
	"strings"

	edgedns "github.com/smart-edge-open/edge-agents/edgedns/pkg/edgednssvr"

	"github.com/miekg/dns"
	"github.com/pkg/errors"
	logger "github.com/smart-edge-open/edge-services/common/log"
	bolt "go.etcd.io/bbolt"
)

var log = logger.DefaultLogger.WithField("storage", nil)

// BoltDB implements the Storage interface
type BoltDB struct {
	Filename string
	instance *bolt.DB
}

var _ edgedns.Storage = &BoltDB{}

// rrSet Resource Records representing the values for a given type
type rrSet struct {
	Rrtype  uint16   // dns.Rrtype
	Answers [][]byte // All answers for a query type
}

const (
	// TTL is the default Time To Live in seconds for authoritative responses
	TTL uint32 = 10

	// "."
	dot = byte(46)

	// Master represents a master (Authoritative) record
	Master uint16 = iota
)

var forwarderBkt = []byte("Forwarders")

// DB Buckets
var bkts = map[uint16]map[uint16][]byte{
	Master: {
		dns.TypeA: {65, 68, 68, 82, 52}, // ADDR4
	},
}

// encode data from a struct
func (rrs *rrSet) encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(&rrs)
	if err != nil {
		log.Errf("Encoding error: %s", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// decode the record into a struct
func decode(data []byte) (*rrSet, error) {
	var rrs *rrSet
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&rrs)
	if err != nil {
		log.Errf("Decoding error: %s", err)
		return nil, err
	}
	return rrs, nil
}

// Start will open the DB file for IO
func (db *BoltDB) Start() error {
	log.Infof("Starting DB from %s", db.Filename)

	var err error
	db.instance, err = bolt.Open(db.Filename, 0660, nil)
	if err != nil {
		return err
	}

	// Create buckets if they do not exist
	err = db.instance.Batch(func(tx *bolt.Tx) error {
		for _, i := range bkts {
			for _, j := range i {
				_, err = tx.CreateBucketIfNotExists(j)
				log.Infof("[DB][%s] Ready", j)
				if err != nil {
					return fmt.Errorf("bucket initialization error: %s", err)
				}
			}
		}

		// Create Forwarders bucket
		if _, err = tx.CreateBucketIfNotExists(forwarderBkt); err != nil {
			return fmt.Errorf("bucket initialization error: %s", err)
		}

		return nil
	})
	return err
}

// Stop will close the DB file from IO
func (db *BoltDB) Stop() error {
	if db.instance != nil {
		if err := db.instance.Close(); err != nil {
			log.Errf("DB Shutdown error: %s", err)
			return err
		}
		return nil
	}
	return errors.New("DB already stopped")
}

// SetHostRRSet creates a resource record
func (db *BoltDB) SetHostRRSet(rrtype uint16,
	fqdn []byte, addrs [][]byte) error {

	if rrtype != dns.TypeA {
		return fmt.Errorf("invalid resource record type (%s),"+
			"only type A supported", dns.TypeToString[rrtype])
	}

	// Make fully qualified
	if !bytes.HasSuffix(fqdn, []byte{dot}) {
		fqdn = append(fqdn, dot)
	}

	err := db.instance.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bkts[Master][rrtype])
		if b == nil {
			return fmt.Errorf("unable to find bucket for %d %d", Master, rrtype)
		}

		for i, j := range addrs {
			log.Debugf("[DB][%d][%d] %d %s: %s",
				Master, rrtype, i+1, fqdn, net.IP(j).String())
		}

		rrs := &rrSet{
			Rrtype:  rrtype,
			Answers: addrs,
		}

		blob, err := rrs.encode()
		if err == nil {
			err = b.Put(fqdn, blob)
		}
		return err
	})
	return err
}

// DelRRSet removes a RR set for a given FQDN and resource type
func (db *BoltDB) DelRRSet(rrtype uint16, fqdn []byte) error {

	// Make fully qualified
	if !bytes.HasSuffix(fqdn, []byte{dot}) {
		fqdn = append(fqdn, dot)
	}

	if _, ok := bkts[Master][rrtype]; ok {
		if err := db.instance.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(bkts[Master][rrtype])
			if b == nil {
				return fmt.Errorf("unable to find bucket for %d %d", Master, rrtype)
			}
			return b.Delete(fqdn)
		}); err != nil {
			return fmt.Errorf("delete %s: %s", fqdn, err)
		}
		log.Debugf("[DB][%d][%d] Delete %s", Master, rrtype, fqdn)
		return nil
	}
	return fmt.Errorf("invalid query type: %s", dns.TypeToString[rrtype])
}

// GetRRSet returns all resources records for an FQDN and resource type
func (db *BoltDB) GetRRSet(name string, rrtype uint16) (*[]dns.RR, error) {
	// Look for Authoritative Answer
	nameLower := strings.ToLower(name)

	rrs := []dns.RR{}
	ans, err := db.getAuthoritative(nameLower, rrtype)
	if err == nil {
		for _, i := range ans.Answers {
			rr, err := rrForType(name, rrtype, i) // nolint: govet
			if err != nil {
				return nil, err
			}
			rrs = append(rrs, rr)
		}
		return &rrs, nil
	}

	return nil, fmt.Errorf("no records found: %w", err)

}

// GetAllRRSets returns all resource records
func (db *BoltDB) GetAllRRSets() (map[string][][]byte, error) {
	allRRs := make(map[string][][]byte)

	if err := db.instance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bkts[Master][dns.TypeA])
		if b == nil {
			return fmt.Errorf("unable to find bucket for %d %d", Master, dns.TypeA)
		}

		if err := b.ForEach(func(fqdn, v []byte) error {
			rs, err := decode(v)
			if err != nil {
				return fmt.Errorf("failed to decode for %s: %w", fqdn, err)
			}

			log.Debugf("[DB][%d][%d] HIT %s", Master, dns.TypeA, fqdn)

			allRRs[string(fqdn)] = rs.Answers

			return nil
		}); err != nil {
			return fmt.Errorf("error calling ForEach: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error calling View: %w", err)
	}

	return allRRs, nil
}

// getAuthoritative returns authoritative records
func (db *BoltDB) getAuthoritative(name string, rrtype uint16) (*rrSet, error) {
	var v []byte

	fqdn := []byte(name)

	if err := db.instance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bkts[Master][rrtype])
		if b == nil {
			return fmt.Errorf("unable to find bucket for %d %d", Master, rrtype)
		}

		v = b.Get(fqdn)

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error calling View: %w", err)
	}

	if len(v) != 0 {
		rrs, err := decode(v)
		if err != nil {
			return nil, fmt.Errorf("failed to decode for %s: %w", fqdn, err)
		}

		log.Debugf("[DB][%d][%d] HIT %s", Master, rrtype, fqdn)
		return rrs, nil
	}

	return nil, fmt.Errorf("no authoritative records found")
}

func rrForType(name string, rrtype uint16, ans []byte) (dns.RR, error) {
	switch rrtype {
	case dns.TypeA:
		r := new(dns.A)
		r.Hdr = dns.RR_Header{
			Name:   name,
			Rrtype: rrtype,
			Class:  dns.ClassINET,
			Ttl:    TTL,
		}
		r.A = net.IP(ans)
		return r, nil
	}
	return nil, fmt.Errorf("uknown resource for Query type: %d", rrtype)
}

// GetForwarders gets the responder's forwarder configuration.
func (db *BoltDB) GetForwarders() ([][]byte, error) {
	var addrs [][]byte

	if err := db.instance.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(forwarderBkt)
		if b == nil {
			return fmt.Errorf("unable to find bucket for %s", forwarderBkt)
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			addrs = append(addrs, v)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("error during view forwarders tx: %v", err)
	}

	return addrs, nil
}

// SetForwarders sets the responder's forwarder configuration.
func (db *BoltDB) SetForwarders(addrs [][]byte) error {
	// Clear all forwarders in bucket
	if err := db.deleteAllKeys(forwarderBkt); err != nil {
		return fmt.Errorf("unable to delete all forwarders: %v", err)
	}

	// If no forwarders are given, return
	if len(addrs) == 0 {
		return nil // no-op
	}

	// Add forwarders to bucket
	return db.instance.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(forwarderBkt)
		if b == nil {
			return fmt.Errorf("unable to find bucket for %s", forwarderBkt)
		}

		for _, addr := range addrs {
			id, err := b.NextSequence()
			if err != nil {
				return fmt.Errorf("unable to get next sequence id: %v", err)
			}

			if err := b.Put(itob(id), addr); err != nil {
				return fmt.Errorf("unable to put forwarder addr %v: %v", addr, err)
			}
		}

		return nil
	})
}

// deleteAllKeys deletes all keys in a bucket.
func (db *BoltDB) deleteAllKeys(bkt []byte) error {
	return db.instance.Batch(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(bkt); err != nil {
			return fmt.Errorf("unable to delete bucket %s: %v", bkt, err)
		}

		if _, err := tx.CreateBucket(bkt); err != nil {
			return fmt.Errorf("unable to recreate bucket %s: %v", bkt, err)
		}

		return nil
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
