<template>
  <v-container grid-list-md text-xs-left>
    <v-layout row wrap>
      <v-flex xs6>
        <div>
          <v-text-field label="Private key" v-model="privateKeyHex"></v-text-field>
          <v-btn @click="generatePrivateKey">Generate ECDSA p-256 key</v-btn>

          <v-text-field label="Public key X" :value="publicKeyX" readonly></v-text-field>
          <v-text-field label="Public key Y" :value="publicKeyY" readonly></v-text-field>
        </div>
      </v-flex>
      <v-flex xs6>
        <div>
          <v-textarea label="User X509" v-model="userCert" rows="15"></v-textarea>
        </div>
      </v-flex>
      <v-flex xs12>
        <hr style="height: 4px; background-color: black;">
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs6>
        <div>
          <v-btn @click="tbsCsr">Get TBS CSR /ca/tbscsr</v-btn>
          <v-text-field label="TBS CSR hash" v-model="tbsCsrHash"></v-text-field>
        </div>
      </v-flex>
      <v-flex xs6>
        <div>
          <v-textarea label="TBS CSR bytes" v-model="tbsCsrBytes" rows="8"></v-textarea>
        </div>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs12>
        <v-btn @click="signTbsCsr">Sign TBS CSR</v-btn>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-flex xs6>
        <v-text-field label="R" v-model="tbsCsrSignature.r"></v-text-field>
      </v-flex>
      <v-flex xs6>
        <v-text-field label="S" v-model="tbsCsrSignature.s"></v-text-field>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs6>
        <v-text-field label="CA login" v-model="caCreds.login"></v-text-field>
      </v-flex>
      <v-flex xs6>
        <v-text-field label="CA password" v-model="caCreds.password"></v-text-field>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs6>
        <v-btn @click="enrollCsr">Enroll to CA with creds, TBS CSR and signature - /ca/enroll_csr</v-btn>
      </v-flex>
      <v-flex xs6>
        <v-alert :value="enrollSuccess" type="success">Cert recieved, user X509 writed.</v-alert>
      </v-flex>
      <v-flex xs12>
        <hr style="height: 4px; background-color: black;">
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs2>
        <v-text-field label="channel_id" v-model="proposalRequest.channel_id"></v-text-field>
      </v-flex>
      <v-flex xs2>
        <v-text-field label="chaincode_id" v-model="proposalRequest.chaincode_id"></v-text-field>
      </v-flex>
      <v-flex xs2>
        <v-text-field label="msp_id" v-model="proposalRequest.msp_id"></v-text-field>
      </v-flex>
      <v-flex xs2>
        <v-text-field label="fcn" v-model="proposalRequest.fcn"></v-text-field>
      </v-flex>
      <v-flex xs2>
        <v-text-field label="args" v-model="proposalRequest.args"></v-text-field>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-flex xs12>
        <v-btn @click="proposal">Call /tx/proposal</v-btn>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs6>
        <div>
          <v-textarea label="Proposal base64 encoded" v-model="proposalB64" rows="20"></v-textarea>
        </div>
      </v-flex>
      <v-flex xs6>
        <div>
          <v-textarea label="Proposal bytes" v-model="proposalBytes" rows="20"></v-textarea>
        </div>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex md6>
        <div>
          <v-text-field label="Proposal hash" v-model="proposalHash"></v-text-field>
        </div>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs12>
        <v-btn @click="signProposal">Sign proposal</v-btn>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-flex xs6>
        <v-text-field label="R" v-model="proposalSignature.r"></v-text-field>
      </v-flex>
      <v-flex xs6>
        <v-text-field label="S" v-model="proposalSignature.s"></v-text-field>
      </v-flex>
      <v-flex xs12>
        <hr style="height: 4px; background-color: black;">
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs12>
        <div>
          <v-btn @click="queryProposal">Query proposal</v-btn>
          <v-text-field label="Query peer" v-model="queryPeer"></v-text-field>
          <v-text-field label="Query response" v-model="queryResponse"></v-text-field>
        </div>
      </v-flex>
      <v-flex xs12>
        <hr style="height: 4px; background-color: black;">
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs12>
        <div>
          <v-btn @click="broadcastPayload">Get Broadcast payload (proposal+endorsments)</v-btn>
          <v-text-field label="Endorsement peers" v-model="endorsementPeers"></v-text-field>
        </div>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs6>
        <div>
          <v-textarea
            label="Broadcast payload base64 encoded"
            v-model="broadcastPayloadB64"
            rows="20"
          ></v-textarea>
        </div>
      </v-flex>
      <v-flex xs6>
        <div>
          <v-textarea label="Broadcast payload bytes" v-model="broadcastPayloadBytes" rows="20"></v-textarea>
        </div>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-flex md6>
        <div>
          <v-text-field label="Broadcast payload hash" v-model="broadcastPayloadHash"></v-text-field>
        </div>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs12>
        <v-btn @click="signBroadcastPayload">Sign broadcast payload</v-btn>
      </v-flex>
    </v-layout>
    <v-layout row wrap>
      <v-flex xs6>
        <v-text-field label="R" v-model="broadcastPayloadSignature.r"></v-text-field>
      </v-flex>
      <v-flex xs6>
        <v-text-field label="S" v-model="broadcastPayloadSignature.s"></v-text-field>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs12>
        <div>
          <v-btn @click="broadcast">Broadcast payload with signature</v-btn>
          <v-text-field label="Broadcast payload response" v-model="broadcastResponse"></v-text-field>
        </div>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
let EC = require("elliptic");
let ec = new EC.ec("p256");

export default {
  data: () => ({
    tbsCsrHash: "",
    tbsCsrBytes: "",

    tbsCsrSignature: {
      r: "",
      s: ""
    },

    caCreds: {
      login: "UserCa",
      password: ""
    },

    enrollSuccess: false,

    privateKeyHex: "",
    userCert: ``,
    proposalB64: "",
    proposalHash: "",
    proposalSignature: {
      r: "",
      s: ""
    },
    proposalRequest: {
      channel_id: "mychannel",
      chaincode_id: "mycc",
      msp_id: "Org1MSP",
      fcn: "query",
      args: "a"
    },

    queryPeer: "org1/peer0",
    queryResponse: "",

    endorsementPeers: "org1/peer0,org2/peer0",

    broadcastPayloadB64: "",
    broadcastPayloadHash: "",
    broadcastPayloadSignature: {
      r: "",
      s: ""
    },

    broadcastResponse: ""
  }),
  computed: {
    broadcastPayloadBytes: function() {
      return atob(this.broadcastPayloadB64);
    },
    proposalBytes: function() {
      return atob(this.proposalB64);
    },
    publicKeyX: function() {
      return this.privateKeyHex.length
        ? ec
            .keyFromPrivate(this.privateKeyHex, "hex")
            .getPublic()
            .getX()
            .toString(16)
        : "";
    },
    publicKeyY: function() {
      return this.privateKeyHex.length
        ? ec
            .keyFromPrivate(this.privateKeyHex, "hex")
            .getPublic()
            .getY()
            .toString(16)
        : "";
    }
  },
  created: function() {
    this.generatePrivateKey();
  },
  methods: {
    tbsCsr: function() {
      let req = {
        X: this.publicKeyX,
        Y: this.publicKeyY,
        Login: this.caCreds.login
      };

      fetch("http://localhost:8080/ca/tbs-csr", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
      })
        .then(response => response.json())
        .then(responseJson => {
          this.tbsCsrBytes = atob(responseJson.tbs_csr_bytes);
          this.tbsCsrHash = responseJson.tbs_csr_hash;
        });
    },
    signTbsCsr: function() {
      let keyPriv = ec.keyFromPrivate(this.privateKeyHex, "hex");
      let signature = keyPriv.sign(this.tbsCsrHash);

      this.tbsCsrSignature.r = signature.r.toString(16);
      this.tbsCsrSignature.s = signature.s.toString(16);
    },
    enrollCsr: function() {
      let req = {
        login: this.caCreds.login,
        password: this.caCreds.password,
        tbs_csr_bytes: btoa(this.tbsCsrBytes),
        r: this.tbsCsrSignature.r,
        s: this.tbsCsrSignature.s
      };

      fetch("http://localhost:8080/ca/enroll-csr", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
      })
        .then(response => response.json())
        .then(responseJson => {
          let pem = atob(responseJson.user_cert);
          if (pem.includes("BEGIN CERTIFICATE")) {
            this.userCert = pem;
            this.enrollSuccess = true;
            setTimeout(() => {
              this.enrollSuccess = false;
            }, 3000);
          }
        });
    },

    generatePrivateKey: function() {
      let key = ec.genKeyPair();
      this.privateKeyHex = key.getPrivate().toString(16);
    },
    proposal: function() {
      let req = { ...this.proposalRequest };
      req.user_cert = this.userCert;
      req.args = req.args.split(",");

      fetch("http://localhost:8080/tx/proposal", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
      })
        .then(response => response.json())
        .then(responseJson => {
          this.proposalB64 = responseJson.proposal_bytes;
          this.proposalHash = responseJson.proposal_hash;
        });
    },
    signProposal: function() {
      let keyPriv = ec.keyFromPrivate(this.privateKeyHex, "hex");
      let signature = keyPriv.sign(this.proposalHash);

      this.proposalSignature.r = signature.r.toString(16);
      this.proposalSignature.s = signature.s.toString(16);
    },
    queryProposal: function() {
      let req = {
        proposal_bytes: this.proposalB64,
        peer: this.queryPeer,
        ...this.proposalSignature
      };

      fetch("http://localhost:8080/tx/query", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
      })
        .then(response => response.json())
        .then(responseJson => {
          this.queryResponse = JSON.stringify(responseJson);
        });
    },
    broadcastPayload: function() {
      let req = {
        proposal_bytes: this.proposalB64,
        peers: this.endorsementPeers.split(","),
        ...this.proposalSignature
      };

      fetch("http://localhost:8080/tx/broadcast-payload", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
      })
        .then(response => response.json())
        .then(responseJson => {
          this.broadcastPayloadB64 = responseJson.payload_bytes;
          this.broadcastPayloadHash = responseJson.payload_hash;
        });
    },
    signBroadcastPayload: function() {
      let ec = new EC.ec("p256");
      let keyPriv = ec.keyFromPrivate(this.privateKeyHex, "hex");
      let signature = keyPriv.sign(this.broadcastPayloadHash);

      this.broadcastPayloadSignature.r = signature.r.toString(16);
      this.broadcastPayloadSignature.s = signature.s.toString(16);
    },
    broadcast: function() {
      let req = {
        payload_bytes: this.broadcastPayloadB64,
        ...this.broadcastPayloadSignature
      };

      fetch("http://localhost:8080/tx/broadcast", {
        method: "POST",
        headers: {
          "Content-Type": "application/json"
        },
        body: JSON.stringify(req)
      })
        .then(response => response.json())
        .then(responseJson => {
          this.broadcastResponse = JSON.stringify(responseJson);
        });
    }
  }
};
</script>

<style>
</style>
