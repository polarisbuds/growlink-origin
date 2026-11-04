const STORAGE_KEY = 'growlink.workspace_tab';

const WORKSPACE_TABS = ['all', 'thisMonth', 'past', 'planned'] as const;
export type WorkspaceTab = (typeof WORKSPACE_TABS)[number];

function isWorkspaceTab(value: string | null): value is WorkspaceTab {
	return WORKSPACE_TABS.some((tab) => tab === value);
}

function readWorkspaceTab(): WorkspaceTab {
	if (typeof localStorage === 'undefined') return 'all';
	const v = localStorage.getItem(STORAGE_KEY);
	if (!isWorkspaceTab(v)) return 'all';
	return v;
}

class WorkspaceStore {
	workspaceTab = $state<WorkspaceTab>(readWorkspaceTab());

	changeTab(tab: WorkspaceTab) {
		this.workspaceTab = tab;
		if (typeof localStorage === 'undefined') return;
		localStorage.setItem(STORAGE_KEY, tab);
	}
}

export const workspace = new WorkspaceStore();
