# ==== Static config ====

# == RTMP/SRT section ==
listen              1935;
max_connections     1000;
daemon              on;

srs_log_tank        file;
srs_log_file        /var/log/srs/srs.log;
srs_log_level       verbose;
ff_log_dir          /var/log/srs/ffmpeg;

http_server {
    enabled         on;
    listen          127.0.0.1:8080;
    crossdomain     on;
    dir             ./objs/nginx/html;
}

http_api {
    enabled         on;
    listen          1985;
    crossdomain     on;
    raw_api {
        enabled             on;
        allow_reload        on;
        allow_query         on;
        allow_update        on;
    }
}

srt_server {
    enabled on;
    listen 9998;
}

vhost __defaultVhost__ {
    mix_correct on;

    http_hooks {
        enabled     on;
        on_publish  http://127.0.0.1:8887/api/v1/webhook/stream/live;
        on_unpublish  http://127.0.0.1:8887/api/v1/webhook/stream/live;
        on_hls  http://127.0.0.1:8887/api/v1/webhook/stream/hls;
    }
    transcode {
        enabled     on;
        ffmpeg      /usr/bin/ffmpeg;

         engine high {
            enabled         on;
            vfilter {
            }
            vcodec          libx264;
            vfps            25;
            vthreads        2;
	        vwidth          -1;
            vheight         1080;
	        vbitrate        6000;
            vprofile        main;
            vpreset         veryfast;
            vparams {
		        maxrate	    6M;
		        bufsize	    3M;
		        g           50;
		        level	    4;
            }
            acodec          aac; 
            abitrate        192;
            asample_rate    44100;
            achannels       2;
            aparams {
            }
            output          rtmp://127.0.0.1:[port]/high?vhost=high/[stream];
        }

        engine mid {
            enabled         on;
            vfilter {
            }
            vcodec          libx264;
            vfps            25;
            vthreads        2;
	        vwidth          -1;
            vheight         720;
	        vbitrate        2000;
            vprofile        main;
            vpreset         veryfast;
            vparams {
		        maxrate	    2M;
		        bufsize	    1M;
		        g 	        50;
		        level	    4;
            }
            acodec          aac;
            abitrate        192;
            asample_rate    44100;
            achannels       2;
            aparams {
            }
            output          rtmp://127.0.0.1:[port]/mid?vhost=mid/[stream];
        }

        engine low {
            enabled         on;
            vfilter {
            }
            vcodec          libx264;
            vfps            25;
            vthreads        2;
	        vwidth          -1;
            vheight         480;
	        vbitrate        1000;
            vprofile        main;
            vpreset         veryfast;
            vparams {
		        maxrate	     1M;
		        bufsize	     500K;
		        g            50;
		        level	     4;
            }
            acodec          aac;
            abitrate        70;
            asample_rate    44100;
            achannels       2;
            aparams {
            }
            output          rtmp://127.0.0.1:[port]/low?vhost=low/[stream];
         }

    }
}

vhost high {
    enabled on;
    mix_correct on;

    hls {
        enabled         on;
        hls_path        ./objs/nginx/html;
        hls_fragment    2;
        hls_window      36000;
	    hls_wait_keyframe       on;
        hls_m3u8_file   live/[stream]/high.m3u8;
        hls_ts_file     live/[stream]/high-[2006]-[01]-[02]-[15]-[04]-[05]-[duration]-[seq].ts;
        hls_acodec      aac;
        hls_vcodec      h264;
    }

    http_hooks {
        enabled     on;
        on_hls  http://127.0.0.1:8887/api/v1/webhook/stream/hls;
    }
}

vhost mid {
    enabled on;
    mix_correct on;

    hls {
        enabled         on;
        hls_path        ./objs/nginx/html;
        hls_fragment    2;
        hls_window      36000;
	    hls_wait_keyframe       on;
        hls_m3u8_file   live/[stream]/mid.m3u8;
        hls_ts_file     live/[stream]/mid-[2006]-[01]-[02]-[15]-[04]-[05]-[duration]-[seq].ts;
        hls_acodec      aac;
        hls_vcodec      h264;
    }

    http_hooks {
        enabled     on;
        on_hls  http://127.0.0.1:8887/api/v1/webhook/stream/hls;
    }
}

vhost low {
    enabled on;
    mix_correct on;

    hls {
        enabled         on;
        hls_path        ./objs/nginx/html;
        hls_fragment    2;
        hls_window      36000;
	    hls_wait_keyframe       on;
        hls_m3u8_file   live/[stream]/low.m3u8;
        hls_ts_file     live/[stream]/low-[2006]-[01]-[02]-[15]-[04]-[05]-[duration]-[seq].ts;
        hls_acodec      aac;
        hls_vcodec      h264;
    }

    http_hooks {
        enabled     on;
        on_hls  http://127.0.0.1:8887/api/v1/webhook/stream/hls;
    }
}

