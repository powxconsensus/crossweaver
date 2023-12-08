package abi

const DIGIPAY_ABI = `[
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "_digiChainId",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "_chainId",
        "type": "string"
      },
      {
        "internalType": "address[]",
        "name": "_validators",
        "type": "address[]"
      }
    ],
    "stateMutability": "nonpayable",
    "type": "constructor"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "string",
        "name": "srcChainId",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "string",
        "name": "dstChainId",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "nonce",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address[]",
        "name": "tokens",
        "type": "address[]"
      },
      {
        "indexed": false,
        "internalType": "uint256[]",
        "name": "amounts",
        "type": "uint256[]"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "sender",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "recipient",
        "type": "address"
      }
    ],
    "name": "Locked",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "internalType": "uint8",
        "name": "txType",
        "type": "uint8"
      },
      {
        "indexed": false,
        "internalType": "string",
        "name": "srcChainId",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "string",
        "name": "dstChainId",
        "type": "string"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "srcNonce",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "uint256",
        "name": "nonce",
        "type": "uint256"
      },
      {
        "indexed": false,
        "internalType": "address[]",
        "name": "tokens",
        "type": "address[]"
      },
      {
        "indexed": false,
        "internalType": "uint256[]",
        "name": "amounts",
        "type": "uint256[]"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "sender",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "address",
        "name": "recipient",
        "type": "address"
      }
    ],
    "name": "UnLocked",
    "type": "event"
  },
  {
    "inputs": [],
    "name": "ETH_ADDRESS",
    "outputs": [
      {
        "internalType": "address",
        "name": "",
        "type": "address"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "uint256",
        "name": "",
        "type": "uint256"
      }
    ],
    "name": "executed",
    "outputs": [
      {
        "internalType": "bool",
        "name": "",
        "type": "bool"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "string",
        "name": "srcChainId",
        "type": "string"
      },
      {
        "internalType": "string",
        "name": "dstChainId",
        "type": "string"
      },
      {
        "internalType": "uint256",
        "name": "srcNonce",
        "type": "uint256"
      },
      {
        "internalType": "bytes",
        "name": "payload",
        "type": "bytes"
      },
      {
        "internalType": "string[]",
        "name": "sigs",
        "type": "string[]"
      }
    ],
    "name": "handleRequest",
    "outputs": [],
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "internalType": "address[]",
        "name": "tokens",
        "type": "address[]"
      },
      {
        "internalType": "uint256[]",
        "name": "amount",
        "type": "uint256[]"
      },
      {
        "internalType": "address",
        "name": "recipient",
        "type": "address"
      }
    ],
    "name": "lock",
    "outputs": [],
    "stateMutability": "payable",
    "type": "function"
  }
]`
