package templates

// OPcacheConfig is the configuration for PHP OPcache
// OPcache improves PHP performance by storing precompiled script bytecode in shared memory
const OPcacheConfig = `
opcache.enable=1
opcache.memory_consumption=256
opcache.interned_strings_buffer=8
opcache.max_accelerated_files=4000
opcache.revalidate_freq=2
opcache.fast_shutdown=1
opcache.enable_cli=1
opcache.validate_timestamps=0
`