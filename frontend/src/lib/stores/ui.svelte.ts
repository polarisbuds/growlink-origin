const STORAGE_KEY = 'growlink.sidebar_open';

function readSidebar(): boolean {
	if (typeof localStorage === 'undefined') return true;
	const v = localStorage.getItem(STORAGE_KEY);
	if (v === 'false') return false;
	if (v === 'true') return true;
	return true;
}

function isMobileViewport(): boolean {
	if (typeof window === 'undefined') return false;
	return window.matchMedia('(max-width: 767px)').matches;
}

class UiStore {
	sidebarOpen = $state(readSidebar());
	mobileDrawerOpen = $state(false);

	toggleSidebar() {
		if (isMobileViewport()) {
			this.mobileDrawerOpen = !this.mobileDrawerOpen;
			return;
		}
		this.sidebarOpen = !this.sidebarOpen;
		localStorage.setItem(STORAGE_KEY, String(this.sidebarOpen));
	}

	closeMobileDrawer() {
		this.mobileDrawerOpen = false;
	}
}

export const ui = new UiStore();
