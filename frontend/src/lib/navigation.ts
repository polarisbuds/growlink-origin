import type { AvatarMe, NavItem } from './types';

export function buildNavigation(pathname: string, me: AvatarMe): NavItem[] {
	const planningActive = pathname.startsWith('/planning');
	const nav: NavItem[] = [
		{
			href: '/',
			label: 'ワークスペース',
			icon: 'home',
			active: pathname === '/' || pathname === '/workspace'
		},
		{
			href: '/planning/meetings',
			label: 'プランニング',
			icon: 'calendar-days',
			active: planningActive
		}
	];
	if (me.isSystemAdmin) {
		nav.push({
			href: '/management',
			label: '管理',
			icon: 'shield',
			active: pathname.startsWith('/management')
		});
	}
	if (me.isSystemAdmin || me.isAnalyst) {
		nav.push({
			href: '/progress',
			label: '進捗',
			icon: 'bar-chart-2',
			active: pathname.startsWith('/progress')
		});
	}
	return nav;
}

export function settingsNav(pathname: string): NavItem {
	return {
		href: '/setting',
		label: '設定',
		icon: 'settings',
		active: pathname.startsWith('/setting')
	};
}
