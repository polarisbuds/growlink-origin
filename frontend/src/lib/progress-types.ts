export type ProgressBreakdownRow = {
	id: number;
	name: string;
	departmentName?: string;
	pending: number;
	draft: number;
	done: number;
	donePercent?: number;
	usersPending?: string[];
	usersDraft?: string[];
	usersDone?: string[];
};
