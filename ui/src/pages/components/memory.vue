<template>
  <div class="card">
    <div class="card-header">
      Memory
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="total_system_memory_human" class="col-sm-7 col-form-label">Total System Memory</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="total_system_memory_human" :value="memory.total_system_memory_human">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="used_memory_rss_human" class="col-sm-7 col-form-label">Used Memory Rss</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="used_memory_rss_human" :value="memory.used_memory_rss_human">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="used_memory_peak_human" class="col-sm-7 col-form-label">Used Memory Peak</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="used_memory_peak_human" :value="memory.used_memory_peak_human">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="used_memory_human" class="col-sm-7 col-form-label">Used Memory</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="used_memory_human" :value="memory.used_memory_human">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="used_memory_lua_human" class="col-sm-7 col-form-label">Used Memory Lua</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="used_memory_lua_human" :value="memory.used_memory_lua_human">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="maxmemory_human" class="col-sm-7 col-form-label">Maxmemory</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="maxmemory_human" :value="memory.maxmemory_human">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="maxmemory_policy" class="col-sm-7 col-form-label">Maxmemory Policy</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="maxmemory_policy" :value="memory.maxmemory_policy">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="mem_fragmentation_ratio" class="col-sm-7 col-form-label">Memory Rragmentation Ratio</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="mem_fragmentation_ratio" :value="memory.mem_fragmentation_ratio">
      </div>
    </div>
    <div class="mb-3 row" style="padding: 0 1.5rem">
      <label for="mem_allocator" class="col-sm-7 col-form-label">Memory Allocator</label>
      <div class="col-sm-5">
        <input type="text" readonly class="form-control-plaintext" id="mem_allocator" :value="memory.mem_allocator">
      </div>
    </div>

  </div>
</template>
<script setup>
import {defineProps} from "vue";

const props = defineProps({
  memory:{}
})

</script>