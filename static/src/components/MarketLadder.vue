<template>
    <div>
        <el-table
            :data="ladderData"
            border>
            <el-table-column
                prop="bidSize"
                label="Bid Size"
                width="170"
                header-align="center"
                align="right">
            </el-table-column>
            <el-table-column
                prop="price"
                label="Price"
                width="159"
                header-align="center">
            </el-table-column>
            <el-table-column
                prop="askSize"
                label="Ask Size"
                width="170"
                header-align="center"
                align="left">
            </el-table-column>
        </el-table>
    </div>
</template>

<script>
export default {
    props: ['ladderData']
}
</script>

<style>
</style>