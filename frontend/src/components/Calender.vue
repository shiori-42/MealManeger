<template>
  <table class="table">
    <tbody>
      <!-- 週ごとの行を動的に生成 -->
      <tr
        v-for="(week, index) in filteredWeeks"
        :key="index"
        :class="`week-${index}`"
      >
        <td class="border-0">
          <!-- 週の最初の日が月の最初の日であれば月名を表示 -->
          <div v-if="isFirstDayOfMonth(week[0].date)" class="fw-bold">
            {{ getMonthName(week[0].date) }}
          </div>
        </td>
        <!-- 日ごとのセルを動的に生成 -->
        <td v-for="day in week" :key="day.date" class="border-0">
          <div
            class="w-100 h-100 d-flex align-items-center justify-content-center"
          >
            <div
              class="border rounded text-center"
              :style="{
                backgroundColor: calculateColor(
                  calculateTotalCalories(day.date.toISOString().slice(0, 10))
                ),
              }"
              style="height: 30px; width: 30px; color: rgb(200, 200, 200)"
            >
              {{ day.date.getDate() }}
            </div>
          </div>
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script>
// var holiday = require("holiday-jp");

export default {
  name: "HelloWorld",
  props: ["userid"],
  data() {
    return {
      weeks: this.generateWeeks(),
      holidays: ["2023-01-01", "2023-12-25"],
      userData: [],
    };
  },
  computed: {
    filteredWeeks() {
      return this.weeks.filter((week) => this.isCase(week[0].date, 12));
    },
  },
  methods: {
    generateWeeks() {
      const year = new Date().getFullYear(); // 現在の年
      let weeks = [];
      let currentWeek = [];
      let currentDate = new Date(year, 0, 1); // 年の最初の日

      while (currentDate.getFullYear() === year) {
        currentWeek.push({
          date: new Date(currentDate),
        });

        if (currentDate.getDay() === 6) {
          // 土曜日で週の終わり
          weeks.push(currentWeek);
          currentWeek = [];
        }

        currentDate.setDate(currentDate.getDate() + 1); // 次の日に移動
      }

      // 最後の週を追加（土曜日で終わらない場合）
      if (currentWeek.length > 0) {
        weeks.push(currentWeek);
      }
      return weeks;
    },
    getMonthName(date) {
      const newDate = new Date(date);
      newDate.setDate(date.getDate() + 6); // 現在の日付から6日後
      return newDate.toLocaleString("en-US", { month: "short" });
    },
    isFirstDayOfMonth(date) {
      const newDate = new Date(date);
      newDate.setDate(date.getDate() + 6); // 現在の日付から6日後
      if (newDate.getFullYear() > 2023) {
        return false;
      }
      return newDate.getDate() <= 7;
    },
    calculateColor(value) {
      if (value == 0) {
        return "#ebedf0";
      }
      if (value <= 2000) {
        return this.interpolateColor("#ffffff", "#fb923c", value / 2000);
      } else if (value <= 3000) {
        return this.interpolateColor(
          "#f87171",
          "#991b1b",
          (value - 2000) / 1000
        );
      } else {
        return "#991b1b";
      }
    },
    interpolateColor(color1, color2, factor) {
      // color1 と color2 間で線形補間を行う
      var result = "#";
      for (var i = 1; i <= 5; i += 2) {
        var color1Part = parseInt(color1.substr(i, 2), 16);
        var color2Part = parseInt(color2.substr(i, 2), 16);
        var mix = Math.round(color1Part + (color2Part - color1Part) * factor);
        result += ("0" + mix.toString(16)).slice(-2);
      }
      return result;
    },
    isCase(date1, mon) {
      const date2 = new Date(date1);
      date2.setDate(date2.getDate() + 6);
      if (mon === 0) {
        // monが0の場合は常にtrue
        return true;
      } else {
        // monが1-12の場合、月をチェック
        const month1 = date1.getMonth() + 1; // JavaScriptの月は0から始まるため、1を加算
        const month2 = date2.getMonth() + 1;
        if (month1 == 12 && month2 == 1) {
          return true;
        }
        return month1 === mon || month2 === mon;
      }
    },
    calculateTotalCalories(date) {
      let totalCalories = 0;
      this.userData.forEach((entry) => {
        if (entry.Date === date) {
          totalCalories += entry.calories;
        }
      });
      return totalCalories;
    },
    fetchUserData() {
      const userid = this.userid;
      console.log("userid=", userid);
      fetch(`http://https://meal-manager-fzgjmrdeka-an.a.run.app/get?userid=${userid}`)
        .then((response) => {
          if (!response.ok) {
            throw new Error("ネットワークレスポンスが不正です。");
          }
          return response.json();
        })
        .then((data) => {
          this.userData = data;
          console.log(data);
        })
        .catch((error) => {
          console.error("データの取得中にエラーが発生しました:", error);
        });
    },
  },
  watch: {
    userid(newVal) {
      if (newVal) {
        this.fetchUserData();
      }
    },
  },
  mounted() {
    if (this.userid) {
      this.fetchUserData();
    }
  },
};
</script>
