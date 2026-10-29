<script lang="ts">
    let { class: className = "h-8 w-8 rounded-full" }: { class?: string } =
        $props();
</script>

<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" class={className}>
    <defs>
        <linearGradient
            id="bgcolor"
            x1="0"
            y1="0"
            x2="512"
            y2="512"
            gradientUnits="userSpaceOnUse"
        >
            <stop offset="0%" style="stop-color: oklch(var(--p));" />
            <stop offset="100%" style="stop-color: oklch(var(--n));" />
        </linearGradient>

        <!-- 温かみのあるインフィニティグラデーション（オレンジ -> ゴールド -> グリーン） -->
        <linearGradient
            id="warmRibbon"
            x1="100"
            y1="255"
            x2="412"
            y2="255"
            gradientUnits="userSpaceOnUse"
        >
            <stop offset="0%" stop-color="#F06A2B" />
            <stop offset="35%" stop-color="#F79B36" />
            <stop offset="55%" stop-color="#F9C449" />
            <stop offset="75%" stop-color="#7BC446" />
            <stop offset="100%" stop-color="#439E54" />
        </linearGradient>

        <!-- リボン内側のハイライト -->
        <linearGradient
            id="warmHighlight"
            x1="100"
            y1="225"
            x2="412"
            y2="285"
            gradientUnits="userSpaceOnUse"
        >
            <stop offset="0%" stop-color="#FF9852" />
            <stop offset="45%" stop-color="#FFE07D" />
            <stop offset="100%" stop-color="#A6E26B" />
        </linearGradient>

        <!-- 交差部分のシャドウ -->
        <linearGradient
            id="underShadow2"
            x1="210"
            y1="215"
            x2="302"
            y2="295"
            gradientUnits="userSpaceOnUse"
        >
            <stop offset="0%" stop-color="#050E1C" stop-opacity="0" />
            <stop offset="32%" stop-color="#050E1C" stop-opacity="0.4" />
            <stop offset="50%" stop-color="#050E1C" stop-opacity="0.52" />
            <stop offset="68%" stop-color="#050E1C" stop-opacity="0.4" />
            <stop offset="100%" stop-color="#050E1C" stop-opacity="0" />
        </linearGradient>

        <clipPath id="centerClip2">
            <rect x="180" y="180" width="152" height="150" />
        </clipPath>

        <!-- 若葉の背後の温かい光彩 -->
        <radialGradient
            id="warmGlow"
            cx="256"
            cy="140"
            r="90"
            gradientUnits="userSpaceOnUse"
        >
            <stop offset="0%" stop-color="#FFB84D" stop-opacity="0.42" />
            <stop offset="60%" stop-color="#FFB84D" stop-opacity="0.12" />
            <stop offset="100%" stop-color="#FFB84D" stop-opacity="0" />
        </radialGradient>
    </defs>

    <!-- 角丸スクエア背景（濃紺） -->
    <rect
        x="-64"
        y="-64"
        width="640"
        height="640"
        rx="143"
        ry="143"
        fill="url(#bgcolor)"
    />

    <!-- ロゴシンボル全体（視覚的な重心に合わせて中央配置） -->
    <g transform="translate(0, 26)">
        <!-- 若葉の背後の光彩 -->
        <circle cx="256" cy="140" r="90" fill="url(#warmGlow)" />

        <!-- インフィニティループ（ベース） -->
        <path
            d="M 256,255
               C 212,195 112,190 112,255
               C 112,320 212,315 256,255
               C 300,195 400,190 400,255
               C 400,320 300,315 256,255 Z"
            fill="none"
            stroke="url(#warmRibbon)"
            stroke-width="38"
            stroke-linecap="round"
            stroke-linejoin="round"
        />

        <path
            d="M 256,255
               C 212,195 112,190 112,255
               C 112,320 212,315 256,255
               C 300,195 400,190 400,255
               C 400,320 300,315 256,255 Z"
            fill="none"
            stroke="url(#warmHighlight)"
            stroke-width="14"
            stroke-linecap="round"
            stroke-linejoin="round"
            opacity="0.65"
        />

        <!-- 下を通るラインのシャドウ -->
        <path
            d="M 256,255 C 212,195 112,190 112,255"
            fill="none"
            stroke="url(#underShadow2)"
            stroke-width="38"
            stroke-linecap="round"
        />
        <path
            d="M 400,255 C 400,320 300,315 256,255"
            fill="none"
            stroke="url(#underShadow2)"
            stroke-width="38"
            stroke-linecap="round"
        />

        <!-- 上を通るライン（立体交差） -->
        <g clip-path="url(#centerClip2)">
            <path
                d="M 112,255
                 C 112,320 212,315 256,255
                 C 300,195 400,190 400,255"
                fill="none"
                stroke="url(#warmRibbon)"
                stroke-width="38"
                stroke-linecap="round"
            />
            <path
                d="M 112,255
                 C 112,320 212,315 256,255
                 C 300,195 400,190 400,255"
                fill="none"
                stroke="url(#warmHighlight)"
                stroke-width="14"
                stroke-linecap="round"
                opacity="0.65"
            />
        </g>

        <!-- 茎と3枚の若葉 -->
        <path
            d="M 256,206 C 256,180 256,158 256,130"
            fill="none"
            stroke="#F8A541"
            stroke-width="8"
            stroke-linecap="round"
        />
        <path
            d="M 256,190 C 244,175 232,163 214,155"
            fill="none"
            stroke="#F69338"
            stroke-width="7"
            stroke-linecap="round"
        />
        <path
            d="M 256,185 C 268,171 282,159 298,151"
            fill="none"
            stroke="#8AD14F"
            stroke-width="7"
            stroke-linecap="round"
        />

        <!-- 左の葉（オレンジゴールド） -->
        <path
            d="M 245,177 C 235,147 212,137 184,141 C 194,167 216,181 245,177 Z"
            fill="#F89B38"
        />
        <path
            d="M 242,175 C 226,163 212,155 198,149"
            fill="none"
            stroke="#FFD27D"
            stroke-width="3"
            stroke-linecap="round"
        />

        <!-- 中央の葉（ゴールデンイエロー） -->
        <path
            d="M 256,151 C 238,123 242,97 256,77 C 270,97 274,123 256,151 Z"
            fill="#FCD053"
        />
        <path
            d="M 256,145 L 256,93"
            fill="none"
            stroke="#FFF0A6"
            stroke-width="3"
            stroke-linecap="round"
        />

        <!-- 右の葉（フレッシュグリーン） -->
        <path
            d="M 267,173 C 276,139 302,127 332,131 C 322,161 298,175 267,173 Z"
            fill="#84CC49"
        />
        <path
            d="M 270,170 C 286,159 302,149 316,141"
            fill="none"
            stroke="#C2F28D"
            stroke-width="3"
            stroke-linecap="round"
        />
    </g>
</svg>
