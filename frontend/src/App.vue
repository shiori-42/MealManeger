<template>
  <div>
    <CalenHeader />
    <Calender :userid="userid" />
  </div>
</template>

<script>
import CalenHeader from './components/CalenHeader.vue';
import Calender from './components/Calender.vue';

export default {
  components: {
    CalenHeader,
    Calender
  },
  data() {
    return {
      userid: null // 初期値をnullに設定
    };
  },
  watch: {
    '$route.query.userid'(newVal) {
      this.userid = newVal;
    }
  },
  mounted() {
    this.userid = this.$route.query.userid; // コンポーネントがマウントされた時にクエリパラメータを読み込む
  }
};
</script>
