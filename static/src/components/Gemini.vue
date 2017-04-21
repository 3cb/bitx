<template>
    <div class="container">
        <h4>This is Gemini Vue</h4>
        <el-button :plain="true" type="success" size="small" :disabled="connected" @click="connect">Connect</el-button>
        <el-button :plain="true" type="danger" size="small" :disabled="!connected" @click="disconnect">Disconnect</el-button>
        <br>
        <el-slider v-model="marketDepth" :min="5" :max="25"></el-slider>
        <br>
        <div class="ml">
            <market-ladder :ladderData="ladderData"></market-ladder>
        </div>
        <time-sales :tradeData="resizedTradeData" class="ts"></time-sales>
        <br>
    </div>
</template>

<script>
import MarketLadder from './MarketLadder.vue'
import TimeSales from './TimeSales.vue'
import xs from 'xstream'
import _ from 'lodash'

export default {
    data() {
        return {
            marketDepth: 7,
            orderBook: [],
            tradeData: [],
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
                        this.$notify({
                            title: "Connected",
                            message: "Streaming data from Gemini Exchange",
                            type: "success"
                        })
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
                        this.$notify({
                            title: "Disconnected",
                            message: "Live data from Gemini Exchange stopped",
                            type: "error"
                        })
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
                            this.orderBook.unshift(x)
                        } else {
                            x.askSize = v.remaining
                            this.orderBook.unshift(x)
                        }
                    })
                    this.orderBook = _.orderBy(this.orderBook, [(o) => parseFloat(o.price)], ["desc"])
                },
                complete: () => {
                    console.log('Order book initialization complete.')
                }
            },
            changeListener: {
                next: (value) => {
                    var i = _.findIndex(this.orderBook, (o) => o.price === value[0].price)

                    switch (i >= 0) {
                        case true:
                            if (value[0].side === "bid") {
                                this.orderBook[i].askSize = ""
                                value[0].remaining === "0" ? _.pullAt(this.orderBook, i) : this.orderBook[i].bidSize = value[0].remaining
                            } else if (value[0].side === "ask") {
                                this.orderBook[i].bidSize = ""
                                value[0].remaining === "0" ? _.pullAt(this.orderBook, i) : this.orderBook[i].askSize = value[0].remaining
                            }
                            break
                        case false:
                            if (value[0].side === "bid") {
                                this.orderBook = _.chain(this.orderBook)
                                                    .concat({
                                                        bidSize: value[0].remaining,
                                                        askSize: "",
                                                        price: value[0].price
                                                    })
                                                    .orderBy([(o) => parseFloat(o.price)], ["desc"])
                                                    .value()
                            } else if (value[0].side === "ask") {
                                this.orderBook = _.chain(this.orderBook)
                                                    .concat({
                                                        bidSize: "",
                                                        askSize: value[0].remaining,
                                                        price: value[0].price
                                                    })
                                                    .orderBy([(o) => parseFloat(o.price)], ["desc"])
                                                    .value()
                            }
                            break
                    }
                },
                complete: () => {
                    console.log('Order Book stream complete.')
                }
            },
            tradeListener: {
                next: (value) => {
                    var x = {
                        amount: value.amount,
                        price: value.price,
                        tid: value.tid
                    }
                    this.tradeData.unshift(x)
                },
                complete: () => {
                    console.log("Time and Sales stream complete.")
                }
            }
        }
    },
    computed: {
        main$() {
            return xs.createWithMemory(this.producerBTC)
        },
        init$() {
            return xs.from(this.main$)
                .take(1)
                .filter(v => v.events.length > 1 && v.type === "update")
                .map(v => v.events)
        },
        change$() {
            return xs.from(this.main$)
                .drop(1)
                .filter(v => v.type === "update")
                .map(v => {
                    if (v.events.length === 2) {
                        _.reverse(_.cloneDeep(v.events))
                    }
                    return v
                })
                .map(v => v.events)
        },
        trade$() {
            return xs.from(this.main$)
                .drop(1)
                .filter(v => v.type === "update")
                .filter(v => v.events.length === 2)
                .map(v => {
                    return v.events[0]
                })
        },
        ladderData() {
            var i = _.findLastIndex(this.orderBook, { "bidSize": "" })
            var j = _.findIndex(this.orderBook, { "askSize": "" })
            return _.slice(this.orderBook, j-this.marketDepth, i+this.marketDepth+1)
        },
        resizedTradeData() {
            return this.tradeData < 2*this.marketData ? this.tradeData: _.dropRight(this.tradeData, this.tradeData.length - 2*this.marketDepth)
        }
    },
    methods: {
        connect() {
            this.connectBTC()
        },
        connectBTC() {
            this.orderBook = []
            this.tradeData = []
            this.main$.addListener(this.controlListener)
            this.init$.addListener(this.initListener)
            this.change$.addListener(this.changeListener)
            this.trade$.addListener(this.tradeListener)
        },
        connectETH() {

        },
        disconnect() {
            this.disconnectBTC()
        },
        disconnectBTC() {
            this.main$.removeListener(this.controlListener)
            this.init$.removeListener(this.initListener)
            this.change$.removeListener(this.changeListener)
            this.trade$.removeListener(this.tradeListener)
        },
        disconnectETH() {

        }
    },
    components: {
        MarketLadder,
        TimeSales
    }
}
</script>

<style>
    .container {
        margin: auto;
        width: 910px;
        align: center;
    }

    .input {
        margin: 0px 0px 0px 0px;
        padding: 0px 0px 0px 0px;
    }

    .ml {
        float: left;
        width: 540px;
    }

    .ts {
        float: right;
    }
</style>