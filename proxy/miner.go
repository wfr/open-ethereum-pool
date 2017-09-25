package proxy

import (
	"log"
	"math/big"
	"strconv"
	"strings"
	//"encoding/hex"
	_"github.com/davecgh/go-spew/spew"

	"github.com/wfr/ethash-nh"
	"github.com/ethereum/go-ethereum/common"

	"github.com/sammy007/open-ethereum-pool/util"
)

var hasher = ethash.New()


// processShare params:
// params[0] = nonce
// params[1] = hashNoNonce
// params[2] = mixDigest

//// NiceHash share processing
func (s *ProxyServer) processShareNH(login, id, ip string, t *BlockTemplate, params[]string) (bool, bool) { // (exist, validShare)
	nonceHex := params[0]
	nonce, _ := strconv.ParseUint(strings.Replace(nonceHex, "0x", "", -1), 16, 64)
	hashNoNonce := common.HexToHash(params[2])

	// Block "difficulty" is BigInt
	// NiceHash "difficulty" is float64 ...
	// diffFloat => target; then: diffInt = 2^256 / target

	shareDiffFloat, mixDigest := hasher.GetShareDiff(t.Height, hashNoNonce, nonce)
	
	// temporary 
	if shareDiffFloat < 0.0001 {
		log.Printf("share difficulty too low, %f < %d, from %v@%v", shareDiffFloat, t.Difficulty, login, ip)
		return false, false
	}
	// temporary hack, ignore round errors
	shareDiffFloat = shareDiffFloat * 0.98

	shareDiff_big := util.DiffFloatToDiffInt(shareDiffFloat)
	shareDiff := shareDiff_big.Int64()

	submit_params := []string{
		nonceHex,
		hashNoNonce.Hex(),
		mixDigest.Hex(),
	}

	h, ok := t.headers[hashNoNonce.Hex()]
	if !ok {
		log.Printf("Stale share from %v@%v", login, ip)
		return false, false
	}

	share := Block{
		number:      h.height,
		hashNoNonce: common.HexToHash(hashNoNonce.Hex()),
		difficulty:  big.NewInt(shareDiff),
		nonce:       nonce,
		mixDigest:   common.HexToHash(mixDigest.Hex()),
	}

	block := Block{
		number:      h.height,
		hashNoNonce: common.HexToHash(hashNoNonce.Hex()),
		difficulty:  h.diff,
		nonce:       nonce,
		mixDigest:   common.HexToHash(mixDigest.Hex()),
	}
	
	if !hasher.Verify(share) {
		log.Println("!hasher.Verify(share)")
		return false, false
	}

	if hasher.Verify(block) {
		ok, err := s.rpc().SubmitBlock(submit_params)
		if err != nil {
			log.Printf("Block submission failure at height %v for %v: %v", h.height, t.Header, err)
		} else if !ok {
			log.Printf("Block rejected at height %v for %v", h.height, t.Header)
			return false, false
		} else {
			s.fetchBlockTemplate()
			exist, err := s.backend.WriteBlock(login, id, submit_params, shareDiff, h.diff.Int64(), h.height, s.hashrateExpiration)
			if exist {
				return true, false
			}
			if err != nil {
				log.Println("Failed to insert block candidate into backend:", err)
			} else {
				log.Printf("Inserted block %v to backend", h.height)
			}
			log.Printf("Block found by miner %v@%v at height %d", login, ip, h.height)
		}
	} else {
		exist, err := s.backend.WriteShare(login, id, submit_params, shareDiff, h.height, s.hashrateExpiration)
		if exist {
			return true, false
		}
		if err != nil {
			log.Println("Failed to insert share data into backend:", err)
		}
	}
	return false, true

}

//// Regular Stratum share processing
func (s *ProxyServer) processShare(login, id, ip string, t *BlockTemplate, params []string) (bool, bool) { // (exist, validShare)
	nonceHex := params[0]
	hashNoNonce := params[1]
	mixDigest := params[2]
	nonce, _ := strconv.ParseUint(strings.Replace(nonceHex, "0x", "", -1), 16, 64)
	shareDiff := s.config.Proxy.Difficulty
	
	h, ok := t.headers[hashNoNonce]
	if !ok {
		log.Printf("Stale share from %v@%v", login, ip)
		return false, false
	}

	share := Block{
		number:      h.height,
		hashNoNonce: common.HexToHash(hashNoNonce),
		difficulty:  big.NewInt(shareDiff),
		nonce:       nonce,
		mixDigest:   common.HexToHash(mixDigest),
	}

	block := Block{
		number:      h.height,
		hashNoNonce: common.HexToHash(hashNoNonce),
		difficulty:  h.diff,
		nonce:       nonce,
		mixDigest:   common.HexToHash(mixDigest),
	}
	
	if !hasher.Verify(share) {
		return false, false
	}

	if hasher.Verify(block) {
		ok, err := s.rpc().SubmitBlock(params)
		if err != nil {
			log.Printf("Block submission failure at height %v for %v: %v", h.height, t.Header, err)
		} else if !ok {
			log.Printf("Block rejected at height %v for %v", h.height, t.Header)
			return false, false
		} else {
			s.fetchBlockTemplate()
			exist, err := s.backend.WriteBlock(login, id, params, shareDiff, h.diff.Int64(), h.height, s.hashrateExpiration)
			if exist {
				return true, false
			}
			if err != nil {
				log.Println("Failed to insert block candidate into backend:", err)
			} else {
				log.Printf("Inserted block %v to backend", h.height)
			}
			log.Printf("Block found by miner %v@%v at height %d", login, ip, h.height)
		}
	} else {
		exist, err := s.backend.WriteShare(login, id, params, shareDiff, h.height, s.hashrateExpiration)
		if exist {
			return true, false
		}
		if err != nil {
			log.Println("Failed to insert share data into backend:", err)
		}
	}
	return false, true
}
