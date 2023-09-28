<template>

  <div>
    <div class="container-fluid text-center">
      <div class="row align-items-start" style="margin: 15px 0;color:#fff;">
        <div class="col" style="background: #0d6efd;height:120px;padding:15px;">
          <div>Queue Total</div>
          <div style="font-weight: bold">{{queue_total}}</div>
        </div>
        <div class="col" style="background: #198754;height:120px;padding:15px;">
          <div>Cpu Usage</div>
          <div style="font-weight: bold">10</div>
        </div>
        <div class="col" style="background: #dc3545;height:120px;padding:15px;">
          <div>Queue Past 10 Minutes</div>
          <div style="font-weight: bold">10</div>
        </div>
        <div class="col" style="background: #343a40;height:120px;padding:15px;">
          <div>DbSize</div>
          <div style="font-weight: bold">{{db_size}}</div>
        </div>
      </div>
    </div>
    <div class="d-flex justify-content-between">
      <div></div>
      <div style="width:50%">
        <v-chart class="chart" :option="barOption" />
      </div>
      <div style="width:50%">
        <v-chart class="chart" :option="lineOption"/>
      </div>
    </div>

  </div>


</template>

<script setup>
import {ref,reactive,onMounted,toRefs,} from "vue";
import request  from "request";

let data = reactive({
  "queue_total":0,
  "db_size":0
})
function getTotal(){
  return request.get("dashboard");
}
onMounted(async ()=>{
  let total = await getTotal();
  data.queue_total = total.data.queue_total;
  data.db_size = total.data.db_size;

})

const barOption = ref({
  title: {
    text: 'Queue Size',
    left: 'left'
  },
  xAxis: {
    type: 'category',
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      data: [120, 200, 150, 80, 70, 110, 130],
      type: 'bar',
      showBackground: true,
      backgroundStyle: {
        color: 'rgba(180, 180, 180, 0.2)'
      }
    }
  ]
});
const lineOption = ref({
  title: {
    text: 'Tasks Processed'
  },
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    bottom:"10",
    data: [
        {
          name:"succeed",
          lineStyle: {
            color: '#198754'
          },
          itemStyle:{
            color:"#198754"
          }
        },
        {
          name:"failed",
          lineStyle:{
            color:'#dc3545'
          },
          itemStyle: {
            color:'#dc3545'
          }
        }
    ]
  },
  toolbox: {
    feature: {
      //saveAsImage: {}
    }
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      name: 'succeed',
      type: 'line',
      lineStyle: {
        color: '#198754'
      },
      itemStyle:{
        color:"#198754"
      },
      data:[220, 182, 191, 234, 290, 330, 310]
    },
    {
      name: 'failed',
      type: 'line',
      lineStyle:{
        color:'#dc3545'
      },
      itemStyle:{
        color:"#dc3545"
      },
      data:[120, 132, 101, 134, 90, 230, 210]
    }
  ]
});
const {queue_total,db_size} = toRefs(data);
</script>
<style scoped>
.chart{
  width:100%;height:80vh;
}
</style>