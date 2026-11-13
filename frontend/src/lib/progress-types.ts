export type ProgressDoneMonthPoint = {
	year: number;
	month: number;
	label: string;
	done: number;
};

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

export type ProgressDepartmentRow = {
	uuid: string;
	name: string;
	depth: number;
	isTotal: boolean;
	userCount: number;
	pairCount: number;
	notDoneCount: number;
	submittedCount: number;
	submittedPercent: number;
	doneCount: number;
	donePercent: number;
	teams?: ProgressTeamRow[];
};

export type ProgressTeamRow = {
	uuid: string;
	name: string;
	isTotal: boolean;
	userCount: number;
	pairCount: number;
	notDoneCount: number;
	submittedCount: number;
	submittedPercent: number;
	doneCount: number;
	donePercent: number;
	departments?: string[];
};

export type ProgressDepartmentTable = {
	rows: ProgressDepartmentRow[];
};
