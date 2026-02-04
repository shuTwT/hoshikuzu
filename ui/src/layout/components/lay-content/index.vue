<script setup lang="ts">
import LayFooter from "../lay-footer/index.vue"
</script>
<template>
    <section class="app-main">
        <n-scrollbar>
            <div class="grow">
                <router-view></router-view>
                
            </div>
            <LayFooter></LayFooter>
        </n-scrollbar>

    </section>
</template>
<style scoped>
.app-main {
    position: relative;
    width: 100%;
    overflow-x: hidden;
    height: 100vh;
    padding-top: 64px;
    background-color: #f5f7fa;
}

.app-main--vertical {
    display: flex;
    flex-direction: column;
}
:deep() .n-scrollbar-content{
    display: flex;
    flex-direction: column;
}
</style>