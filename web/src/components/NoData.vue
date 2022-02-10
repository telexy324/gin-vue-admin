<template lang="pug">
    //- 暂无数据
    .NoData
        .title {{title}}
            slot
</template>

<script>
export default {
    name: 'NoData',
    props: {
        title: {
            type: String,
            default: 'No Data',
        },
    },
    data() {
        return {}
    },
    computed: {},
    created() {},
    mounted() {},
    watch: {},
    methods: {},
    components: {},
}
</script>

<style scoped lang="stylus">
.NoData {
    text-align center
    margin 100px 0
    img {
        width 100px
    }
    .title {
        font-size 18px
        color #8497a4
    }
}
</style>
