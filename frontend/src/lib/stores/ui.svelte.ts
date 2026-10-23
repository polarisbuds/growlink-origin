const STORAGE_KEY = 'growlink.sidebar_open';

function readSidebar(): boolean {
	if (typeof localStorage === 'undefined') return true;
	const v = localStorage.getItem(STORAGE_KEY);
	if (v === 'false') return false;
	if (v === 'true') return true;
	return true;
}

class UiStore {
	sidebarOpen = $state(readSidebar());

	toggleSidebar() {
		this.sidebarOpen = !this.sidebarOpen;
		localStorage.setItem(STORAGE_KEY, String(this.sidebarOpen));
	}
}

export const ui = new UiStore();
