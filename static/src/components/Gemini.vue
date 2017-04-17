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
                        listener.next(event)
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
                next: (v) => {

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
                    .map(event => {
                        let x = JSON.parse(event.data)
                        console.log(x)
                        return x
                    })
        },
        init$() {
            return xs.from(this.parsed$).take(1)
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

            // change$.addListener(this.changeListener)


            // listen for messages and update order book
            // this.gemSocketBTC.onmessage = (event) => {
            //     var msg = JSON.parse(event.data)
            //     console.log(msg)
            //     var hi = msg.events[msg.events.length - 1].price
            //     var lo = msg.events[0].price



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