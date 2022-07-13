package ethofs

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"

	cid "github.com/ipfs/go-cid"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	options "github.com/ipfs/interface-go-ipfs-core/options"
	path "github.com/ipfs/interface-go-ipfs-core/path"
)

func updateLocalPinMapping(api coreiface.CoreAPI) (error, bool) {

	tempPinMapping := make(map[string]string)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opt, err := options.Pin.Ls.Type("all")
	if err != nil {
                log.Debug("ethoFS - local pin mapping option failure", "error", err)
		return err, false
	}

	pins, err := api.Pin().Ls(ctx, opt)
	if err != nil {
                log.Debug("ethoFS - local pin mapping failure", "error", err)
		return err, false
	}


	tempCount := 0

	for p := range pins {
		if p != nil && p.Path() != nil {
			tempCount++
			tempPinMapping[p.Path().Cid().String()] = p.Type()
		}
	}

	select {
	case <-ctx.Done():
		if len(tempPinMapping) > 0 {
			localPinMapping = tempPinMapping
			log.Info("ethoFS - local pin mapping complete", "pin count", len(localPinMapping))
			return nil, false
		} else {
                	log.Debug("ethoFS - local pin mapping failure", "pin count", len(localPinMapping))
       	        	return fmt.Errorf("Error - ethoFS local pin mapping failure"), false
		}
       	}
	return nil, false
}

func pinSearch(hash string, pinMapping map[string]string) bool {

	tempPinMapping := pinMapping

	if _, found := tempPinMapping[hash]; found {
		log.Debug("ethoFS - Matching pin found", "type", found, "hash", hash)
		return true
	}
	log.Debug("ethoFS - Matching pin was not found", "hash", hash)
	return false
}

func pinAdd(api coreiface.CoreAPI, hash string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cid, err := cid.Parse(hash)
	if err != nil {
		return hash, err
	}

	resolvedPath := path.IpfsPath(cid)

	if err := api.Pin().Add(ctx, resolvedPath, options.Pin.Recursive(true)); err != nil {
		return hash, err
	}

	return hash, nil
}

func pinRemove(api coreiface.CoreAPI, hash string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cid, err := cid.Parse(hash)
	if err != nil {
		return hash, err
	}

	resolvedPath := path.IpfsPath(cid)

	if err := api.Pin().Rm(ctx, resolvedPath, options.Pin.RmRecursive(true)); err != nil {
		return hash, err
	}

	return hash, nil
}
