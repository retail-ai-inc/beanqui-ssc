<template>
  <div class="row align-items-start dashboard-index" style="margin: 1.25rem 0;color:#fff;">
    <div class="col" style="background: #0d6efd">
      <div>Queue Total</div>
      <div style="font-weight: bold">
        <router-link to="/admin/queue" class="nav-link text-muted link-color" >{{queue_total}}</router-link></div>
    </div>
    <div class="col" style="background: #198754">
      <div>CPU Total</div>
      <div style="font-weight: bold">
        <router-link to="/admin/redis" class="nav-link text-muted link-color">{{num_cpu}}</router-link>
      </div>
    </div>
    <div class="col" style="background: #dc3545">
      <div>Fail Total</div>
      <div style="font-weight: bold">
        <router-link to="" class="nav-link text-muted link-color">{{fail_count}}</router-link>
      </div>
    </div>
    <div class="col" style="background: #20c997">
      <div>Success Total</div>
      <div style="font-weight: bold">
        <router-link to="" class="nav-link text-muted link-color">{{success_count}}</router-link>
      </div>
    </div>
    <div class="col" style="background: #343a40">
      <div>Total Payload</div>
      <div style="font-weight: bold">
        <router-link to="" class="nav-link text-muted link-color">{{db_size}}</router-link>
      </div>
    </div>
  </div>
</template>
<script setup>
import {defineProps} from "vue";

const props = defineProps({
  queue_total:0,
  num_cpu:0,
  fail_count:0,
  success_count:0,
  db_size:0
})

</script>
<style scoped>
.dashboard-index .col{
  height:7.5rem;
  padding:1rem;
}
.link-color{
  display: inline-block;
  color: #fff !important;
}
</style>