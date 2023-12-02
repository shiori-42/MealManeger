<template>
  <div id="app">
    <calorie-graph></calorie-graph>
  </div>
</template>

<script>
import CalorieGraph from './components/CalorieGraph.vue';

export default {
  name: 'App',
  components: {
    CalorieGraph
  }
};
</script>
