<template>
    <div>
        <h1>This is Gemini Vue</h1>
        <el-button :plain="true" type="success" :disabled="connected" @click="connect">Connect</el-button>
        <el-button :plain="true" type="danger" :disabled="!connected" @click="disconnect">Disconnect</el-button>
    </div>
</template>

<script>
import xs from 'xstream'

export default {
    data() {
        return {
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
                            bid: "",
                            ask: "",
                            remaining: v.remaining
                        }
                        if (v.side === "bid") {
                            x.bid = v.price
                            this.bidIndex.unshift(parseFloat(v.price))
                            this.bidData.unshift(x)
                        } else {
                            x.ask = v.price
                            this.askIndex.push(parseFloat(v.price))
                            this.askData.push(x)
                        }
                    })
                        console.log(this.bidData)
                        console.log(this.askData)
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



                // if (msg.type === "update" && msg.events.length > 1) {
                //     msg.events.forEach((v) => {
                //         let x = {
                //             price: v.price,
                //             remaining: v.remaining
                //         }
                //         if (v.side === "bid") {
                //             this.bidIndex.unshift(parseInt(x.price))
                //             this.bidData.unshift(x)
                //         } else {
                //             this.askIndex.push(parseInt(x.price))
                //             this.askData.push(x)
                //         }
                //     })
                // } else {}
            // }
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
    }
}
</script>