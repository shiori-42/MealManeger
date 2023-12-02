<template>
  <div class="calorie-graph">
    <div v-for="week in calculatedWeeks" :key="week.id" class="week">
      <div
        v-for="day in week.days"
        :key="day.date"
        class="day"
        :style="{ backgroundColor: getDayColor(day.value) }"
      ></div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      // ここではダミーデータを使用しています
      dailyCalories: this.generateDummyData(),
    };
  },
  computed: {
    calculatedWeeks() {
      let weeks = [];
      for (let i = 0; i < this.dailyCalories.length; i += 7) {
        weeks.push({
          id: i,
          days: this.dailyCalories.slice(i, i + 7),
        });
      }
      return weeks;
    },
  },
  methods: {
    generateDummyData() {
      // ダミーデータの生成（365日分のランダムなカロリーデータ）
      return Array.from({ length: 365 }, () => ({
        date: new Date(),
        value: Math.random() * 500,
      }));
    },
    getDayColor(calories) {
      const intensity = Math.min(255, calories / 2); // カロリーに基づいて色の強度を計算
      return `rgba(255, 165, 0, ${intensity / 255})`; // オレンジ色
    },
  },
};
</script>

<style>
.calorie-graph {
  display: flex;
  flex-direction: column;
}
.week {
  display: flex;
}
.day {
  width: 20px;
  height: 20px;
  margin: 2px;
  border: 1px solid #ddd;
}
@media (max-width: 600px) {
  .day {
    width: 15px;
    height: 15px;
  }
}
</style>
