<template>
    <svg width="450" height="450">
        <!--header-->
        <g>
            <rect width="450" height="30"></rect>
            <text text-anchor="end" x="140" y="20">Bid Size</text>
            <text text-anchor="middle" x="225" y="20">Price</text>
            <text text-anchor="start" x="310" y="20">Ask Size</text>
        </g>

        <!--horizontal lines and data-->
        <g>
            <g v-for="(obj, index) in ladderData">
                <text text-anchor="end" v-bind="{ 'x':140, 'y':(index+1)*(420/ladderData.length)-10+30 }">{{ obj.bidSize }}</text>
                <text text-anchor="middle" v-bind="{ 'x':225, 'y':(index+1)*(420/ladderData.length)-10+30 }">{{ obj.price }}</text>
                <text text-anchor="start" v-bind="{ 'x':310, 'y':(index+1)*(420/ladderData.length)-10+30 }">{{ obj.askSize }}</text>
                <line x1="0" x2="450" v-bind="{ 'y1':(index+1)*(420/ladderData.length)+30, 'y2':(index+1)*(420/ladderData.length)+30 }"></line>
            </g>
        </g>

        <!--vertical lines-->
        <g>
            <line x1="0" x2="0" y1="0" y2="450"></line>
            <line x1="150" x2="150" y1="0" y2="450"></line>
            <line x1="300" x2="300" y1="0" y2="450"></line>
            <line x1="450" x2="450" y1="0" y2="450"></line>
        </g>
    </svg>
</template>

<script>
export default {
    props: ['ladderData']
}
</script>
  
<style>
  /*.clearfix:before,
  .clearfix:after {
      display: table;
      content: "";
  }
  .clearfix:after {
      clear: both
  }*/
</style>