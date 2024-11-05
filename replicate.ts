const keys_values = {};
const nodes = {}; // Maybe want a [] instead

// When any key or node change happens, we loop through all nodes
// and send an internal API call to all other nodes?
function add_key(key, value, should_replicate: boolean) {}

function delete_key(key, should_replicate: boolean) {}

function get_key(key, should_replicate: boolean) {}

function add_node(should_replicate: boolean) {}

function delete_node(should_replicate: boolean) {}

function get_node(should_replicate: boolean) {}
