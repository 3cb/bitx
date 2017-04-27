<template>
    <div class="container">
        <h4>Gemini Exchange</h4>
        <br>
        <span>last: {{ initPrice.last }}  </span><span>bid: {{ initPrice.bid }}  </span><span>ask: {{ initPrice.ask }}  </span><span>24-hour volume: {{ initPrice.volume.BTC }}</span>
        <br>
        <el-button :plain="true" type="success" size="small" :disabled="connected" @click="connect">Connect</el-button>
        <el-button :plain="true" type="danger" size="small" :disabled="!connected" @click="disconnect">Disconnect</el-button>
        <el-radio-group v-model="radioControl" @change="switchSocket">
            <el-radio label="btcusd">BTC/USD</el-radio>
            <el-radio label="ethusd">ETH/USD</el-radio>
            <el-radio label="ethbtc">ETH/BTC</el-radio>
        </el-radio-group>
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
import axios from 'axios'
import getTime from 'date-fns/get_time'
import subMonths from 'date-fns/sub_months'
import xs from 'xstream'
import _ from 'lodash'

export default {
    data() {
        return {
            currency: "btcusd",
            radioControl: "btcusd",
            initPrice: {
                bid: "",
                ask: "",
                last: "",
                volume: { BTC: "" }
            },
            tradeHistory: [],
            historyMerged: false,
            marketDepth: 7,
            orderBook: [],
            tradeData: [],
            connected: false,
            autoReconnect: false,
            gemSocket: '',
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
                    if (this.historyMerged === false) {
                        let z = this.tradeData[this.tradeData.length-1].tid
                        this.tradeData = _.concat(this.tradeData, _.dropWhile(this.tradeHistory, o => o.tid >= z))
                        this.historyMerged = true
                    }
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
        websocketAddr() {
            return 'wss://api.gemini.com/v1/marketdata/' + this.currency
        },
        producer() {
            return {
                start: (listener) => {
                    this.gemSocket = new WebSocket(this.websocketAddr)
                    this.gemSocket.onopen = (event) => {
                        this.orderBook = []    
                        // this.tradeData = []
                        this.connected = true
                        this.$notify({
                            title: "Connected",
                            message: "Streaming data from Gemini Exchange",
                            type: "success"
                        })
                        console.log(event)
                    }
                    this.gemSocket.onmessage = (event) => {
                        listener.next(JSON.parse(event.data))
                    }
                },
                stop: () => {
                    this.gemSocket.close()
                    this.gemSocket.onclose = (event) => {
                        this.connected = false
                        this.$notify({
                            title: "Disconnected",
                            message: "Live data from Gemini Exchange stopped",
                            type: "error"
                        })
                        console.log(event)
                    }
                }
            }
        },
        main$() {
            return xs.createWithMemory(this.producer)
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
                        var x = _.cloneDeep(v)
                        _.reverse(x.events)
                        return x
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
                .map(v => v.events[0])
        },
        ladderData() {
            var i = _.findLastIndex(this.orderBook, { "bidSize": "" })
            var j = _.findIndex(this.orderBook, { "askSize": "" })
            return _.slice(this.orderBook, j-this.marketDepth, i+this.marketDepth+1)
        },
        resizedTradeData() {
            return this.tradeData < 2*this.marketData ? this.tradeData : _.dropRight(this.tradeData, this.tradeData.length - 2*this.marketDepth)
        }
    },
    methods: {
        connect() {
            this.tradeData = []
            this.tradeHistory = []
            this.main$.addListener(this.controlListener)
            this.init$.addListener(this.initListener)
            this.change$.addListener(this.changeListener)
            this.trade$.addListener(this.tradeListener)
            this.getHistory()
        },
        disconnect() {
            this.main$.removeListener(this.controlListener)
            this.init$.removeListener(this.initListener)
            this.change$.removeListener(this.changeListener)
            this.trade$.removeListener(this.tradeListener)
            this.historyMerged = false
            this.tradeData = []
            this.tradeHistory = []
        },
        switchSocket() {
            this.currency = this.radioControl
            this.autoReconnect = true
            if (this.connected === true) {
                this.disconnect()
            }
        },
        getHistory() {
            axios.get("/api/getHistory/" + this.currency)
            .then(response => {
                response.data.forEach(obj => {
                    this.tradeHistory.push({amount: obj.amount, price: obj.price, tid: obj.tid})
                })
                if (!this.tradeData.length) {
                    this.tradeData = this.tradeHistory
                    this.historyMerged = true
                }
            })
            .catch(error => {
                console.log("Unable to retrieve historical data: " + error)
            })
        }
    },
    beforeCreate() {
        axios.get("/api/initializePrice")
            .then(response => {
                // console.log(response.data)
                this.initPrice = response.data
            })
            .catch(error => {
                console.log(error)
            })
    },
    beforeUpdate() {
        if (this.autoReconnect === true && this.connected === false) {
            this.autoReconnect = false
            this.connect()
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