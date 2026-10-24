
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	type MatcherParam<M> = M extends (param : string) => param is (infer U extends string) ? U : string;

	export interface AppTypes {
		RouteId(): "/" | "/management" | "/meeting" | "/meeting/new" | "/meeting/[uuid]" | "/meeting/[uuid]/answer" | "/meeting/[uuid]/setting" | "/planning" | "/planning/meetings" | "/planning/meetings/new" | "/planning/pairs" | "/planning/pairs/new" | "/planning/pairs/[uuid]" | "/planning/pairs/[uuid]/edit" | "/planning/templates" | "/planning/templates/new" | "/planning/templates/[uuid]" | "/progress" | "/setting";
		RouteParams(): {
			"/meeting/[uuid]": { uuid: string };
			"/meeting/[uuid]/answer": { uuid: string };
			"/meeting/[uuid]/setting": { uuid: string };
			"/planning/pairs/[uuid]": { uuid: string };
			"/planning/pairs/[uuid]/edit": { uuid: string };
			"/planning/templates/[uuid]": { uuid: string }
		};
		LayoutParams(): {
			"/": { uuid?: string | undefined };
			"/management": Record<string, never>;
			"/meeting": { uuid?: string | undefined };
			"/meeting/new": Record<string, never>;
			"/meeting/[uuid]": { uuid: string };
			"/meeting/[uuid]/answer": { uuid: string };
			"/meeting/[uuid]/setting": { uuid: string };
			"/planning": { uuid?: string | undefined };
			"/planning/meetings": Record<string, never>;
			"/planning/meetings/new": Record<string, never>;
			"/planning/pairs": { uuid?: string | undefined };
			"/planning/pairs/new": Record<string, never>;
			"/planning/pairs/[uuid]": { uuid: string };
			"/planning/pairs/[uuid]/edit": { uuid: string };
			"/planning/templates": { uuid?: string | undefined };
			"/planning/templates/new": Record<string, never>;
			"/planning/templates/[uuid]": { uuid: string };
			"/progress": Record<string, never>;
			"/setting": Record<string, never>
		};
		Pathname(): "/" | "/management" | "/meeting" | "/meeting/new" | `/meeting/${string}` & {} | `/meeting/${string}/answer` & {} | `/meeting/${string}/setting` & {} | "/planning/meetings" | "/planning/meetings/new" | "/planning/pairs" | "/planning/pairs/new" | `/planning/pairs/${string}/edit` & {} | "/planning/templates" | "/planning/templates/new" | `/planning/templates/${string}` & {} | "/progress" | "/setting";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): string & {};
	}
}