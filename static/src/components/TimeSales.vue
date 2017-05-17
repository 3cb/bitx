<template>
    <svg width="360" height="450">
        <!--header-->
        <g>
            <rect width="360" height="30"></rect>
            <text text-anchor="middle" x="90" y="20">Price</text>
            <text text-anchor="middle" x="270" y="20">Amount</text>
        </g>

        <!--horizontal lines and data-->
        <g>
            <g v-for="(obj, index) in tradeData">
                <text text-anchor="middle" v-bind="{ 'x':90, 'y':(index+1)*(420/tradeData.length)+20, 'fill':obj.color }">{{ obj.price }}</text>
                <text text-anchor="middle" v-bind="{ 'x':270, 'y':(index+1)*(420/tradeData.length)+20, 'fill':obj.color }">{{ obj.amount }}</text>
            </g>
            <line x1="0" x2="360" y1="450" y2="450"></line>
        </g>

        <!--vertical lines-->
        <g>
            <line x1="0" x2="0" y1="0" y2="450"></line>
            <line x1="360" x2="360" y1="0" y2="450"></line>
        </g>
    </svg>
</template>

<script>
export default {
    props: ['tradeData']
}
</script>

<style>
    svg {
        background-color: hsl(0, 0%, 96%);
    }
</style>