<template>
  <CalenHeader />
  <Calender />
</template>

<script>
import Calender from './components/Calender.vue'
import CalenHeader from './components/CalenHeader.vue'

export default {
  name: 'App',
  components: {
    Calender,CalenHeader
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
