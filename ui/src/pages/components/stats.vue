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
      <label for="netInput" class="col-sm-7 col-form-label">Total Net Input Bytes</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="netInput" :value="`${stats.total_net_input_bytes}byte`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="netOutput" class="col-sm-7 col-form-label">Total Net Output Bytes</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="netOutput" :value="`${stats.total_net_output_bytes}byte`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="inputKbps" class="col-sm-7 col-form-label">Instantaneous Input Kbps</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="inputKbps" :value="`${stats.instantaneous_input_kbps}`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="outputKbps" class="col-sm-7 col-form-label">Instantaneous Input Kbps</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="outputKbps" :value="`${stats.instantaneous_output_kbps}`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="rejConn" class="col-sm-7 col-form-label">Rejected Connections</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="rejConn" :value="`${stats.rejected_connections}`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="expiredKeys" class="col-sm-7 col-form-label">Expired Keys</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="expiredKeys" :value="`${stats.expired_keys}`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="evictedKeys" class="col-sm-7 col-form-label">Evicted Keys</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="evictedKeys" :value="`${stats.evicted_keys}`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="keyspaceHits" class="col-sm-7 col-form-label">Keyspace Hits</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="keyspaceHits" :value="`${stats.keyspace_hits}`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="keyspaceMisses" class="col-sm-7 col-form-label">Keyspace Misses</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="keyspaceMisses" :value="`${stats.keyspace_misses}`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="latestForkUsec" class="col-sm-7 col-form-label">Latest Fork Usec</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="latestForkUsec" :value="`${stats.latest_fork_usec}`">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="migrateCachedSockets" class="col-sm-7 col-form-label">Migrate Cached Sockets</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="migrateCachedSockets" :value="`${stats.migrate_cached_sockets}`">
      </div>
    </div>

  </div>
</template>
<script setup>
import {defineProps} from "vue";

const props = defineProps({
  stats:{}
})

</script>