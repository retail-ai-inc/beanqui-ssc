<template>
  <div class="card">
    <div class="card-header">
      Stats
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="connReceived" class="col-sm-7 col-form-label">Total Connections Received</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="connReceived" :value="stats.total_connections_received">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="commandProcessed" class="col-sm-7 col-form-label">Total Commands Processed</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="commandProcessed" :value="stats.total_commands_processed">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="instPerSec" class="col-sm-7 col-form-label">Instantaneous Ops Per Sec</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="instPerSec" :value="`${stats.instantaneous_ops_per_sec}/s`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="netInput" class="col-sm-7 col-form-label">total_net_input_bytes</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="netInput" :value="`${stats.total_net_input_bytes}`">
      </div>
    </div>





    <!--            total_net_output_bytes:66358938 #redis网络出口流量字节数-->
    <!--            instantaneous_input_kbps:0.04  #redis网络入口kps-->
    <!--            instantaneous_output_kbps:3633.35  #redis网络出口kps-->
    <!--            rejected_connections:0  #拒绝的连接个数，redis连接个数达到maxclients限制，拒绝新连接的个数-->
    <!--            sync_full:0  #主从完全同步成功次数-->
    <!--            sync_partial_ok:0  #主从部分同步成功次数-->
    <!--            sync_partial_err:0  #主从部分同步失败次数-->
    <!--            expired_keys:0   #运行以来过期的key的数量-->
    <!--            evicted_keys:0  #运行以来剔除(超过了maxmemory后)的key的数量-->
    <!--            keyspace_hits:87  #命中次数-->
    <!--            keyspace_misses:17   #没命中次数-->
    <!--            pubsub_channels:0  #当前使用中的频道数量-->
    <!--            pubsub_patterns:0  #当前使用的模式的数量-->
    <!--            latest_fork_usec:0   #最近一次fork操作阻塞redis进程的耗时数，单位微秒-->
    <!--            migrate_cached_sockets:0   #是否已经缓存了到该地址的连接-->

  </div>
</template>
<script setup>
import {defineProps} from "vue";

const props = defineProps({
  stats:{}
})

</script>