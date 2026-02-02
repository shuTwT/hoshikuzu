<script setup lang="ts">
import LayFooter from "../lay-footer/index.vue"
</script>
<template>
    <section class="app-main">
        <router-view></router-view>
        <LayFooter></LayFooter>
    </section>
</template>
<style scoped>
.app-main {
    position: relative;
    width: 100%;
    height: 100vh;
    overflow-x: hidden;
}
</style>