<template>
    <div class="container">
        <h4>This is Gemini Vue</h4>
        <el-button :plain="true" type="success" size="small" :disabled="connected" @click="connect">Connect</el-button>
        <el-button :plain="true" type="danger" size="small" :disabled="!connected" @click="disconnect">Disconnect</el-button>
        <br>
        <el-slider v-model="marketDepth" :min="5" :max="25"></el-slider>
        <br>
        <market-ladder :ladderData="ladderData"></market-ladder>
        <br>
    </div>
</template>

<script>
import MarketLadder from './MarketLadder.vue'
import xs from 'xstream'
import _ from 'lodash'

export default {
    data() {
        return {
            marketDepth: 10,
            bidIndex: [],
            bidData: [],
            askIndex: [],
            askData: [],
            connected: false,
            geminiAddrBTC: 'wss://api.gemini.com/v1/marketdata/btcusd',
            geminiAddrETH: 'wss://api.gemini.com/v1/marketdata/ethusd',
            gemSocketBTC: '',
            gemSocketETH: '',
            producerBTC: {
                start: (listener) => {
                    this.gemSocketBTC = new WebSocket(this.geminiAddrBTC)
                    this.gemSocketBTC.onopen = (event) => {
                        this.connected = true
                        console.log(event)
                    }
                    this.gemSocketBTC.onmessage = (event) => {
                        listener.next(JSON.parse(event.data))
                    }
                },
                stop: () => {
                    this.gemSocketBTC.close()
                    this.gemSocketBTC.onclose = (event) => {
                        this.connected = false
                        console.log(event)
                    }
                }
            },
            controlListener: {
                next: () => { return },
                error: (err) => {
                    console.error('Error from websocket: ', err)
                },
                complete: () => {
                    console.log('Stream Complete.')
                }
            },
            initListener: {
                next: (value) => {
                    value.forEach(v => {
                        var x = {
                            bidSize: "",
                            askSize: "",
                            price: v.price
                        }
                        if (v.side === "bid") {
                            x.bidSize = v.remaining
                            this.bidIndex.unshift(parseFloat(x.price))
                            this.bidData.unshift(x)
                        } else {
                            x.askSize = v.remaining
                            this.askIndex.unshift(parseFloat(x.price))
                            this.askData.unshift(x)
                        }
                    })
                },
                complete: () => {
                    console.log('Order book initialization complete.')
                }
            },
            changeListener: {
                next: (v) => {

                },
                complete: () => {
                    console.log('Price update stream complete.')
                }
            }
        }
    },
    computed: {
        parsed$() {
            return xs.createWithMemory(this.producerBTC)
                    .map(v => {
                        console.log(v)
                        return v
                    })
        },
        init$() {
            return xs.from(this.parsed$)
                .filter(v => v.type === "update")
                .map(v => v.events)
                .take(1)
        },
        change$() {
            return xs.from(this.parsed$).drop(1)
        },
        ladderData() {
            return _.chain(this.askData)
                        .drop(this.askData.length - this.marketDepth)
                        .concat(_.dropRight(this.bidData, this.bidData.length - this.marketDepth))
                        .value()
        }
        
    },
    methods: {
        connect() {
            this.connectBTC()
        },
        connectBTC() {
            this.parsed$.addListener(this.controlListener)
            this.init$.addListener(this.initListener)
            this.change$.addListener(this.changeListener)
        },
        connectETH() {

        },
        disconnect() {
            this.disconnectBTC()
        },
        disconnectBTC() {
            this.parsed$.removeListener(this.controlListener)
            this.init$.removeListener(this.initListener)
            this.change$.removeListener(this.changeListener)
        },
        disconnectETH() {

        }
    },
    components: {
        MarketLadder
    }
}
</script>

<style>
    .container {
        width: 500px;
    }

    .input {
        margin: 0px 0px 0px 0px;
        padding: 0px 0px 0px 0px;
    }
</style>