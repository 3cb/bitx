<template>
    <div class="ts">
        <el-table
            :data="tradeData"
            maxHeight="600">
            <el-table-column
                prop="price"
                label="Price"
                width="184"
                align="center"
                header-align="center">
            </el-table-column>
            <el-table-column
                prop="amount"
                label="Size"
                width="184"
                align="center"
                header-align="center">
            </el-table-column>
        </el-table>
    </div>
</template>

<script>
export default {
    props: ['tradeData']
}
</script>

<style>
  .clearfix:before,
  .clearfix:after {
      display: table;
      content: "";
  }
  .clearfix:after {
      clear: both
  }
  .ts {
      width: 370px;
  }
</style>