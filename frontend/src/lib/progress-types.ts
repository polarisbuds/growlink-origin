export type ProgressBreakdownRow = {
	uuid: string;
	name: string;
	departmentName?: string;
	pending: number;
	draft: number;
	submitted: number;
	done: number;
	donePercent?: number;
	usersPending?: string[];
	usersDraft?: string[];
	usersSubmitted?: string[];
	usersDone?: string[];
};
