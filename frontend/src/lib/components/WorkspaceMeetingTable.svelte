<script lang="ts">
	import { statusBadge, statusLabel } from '$lib/status';
	import { formatDate, pairLabel, type MeetingListItem } from '$lib/meeting-helpers';

	let {
		items,
		showAnswer = true,
		showSetting = true
	}: {
		items: MeetingListItem[];
		showAnswer?: boolean;
		showSetting?: boolean;
	} = $props();
</script>

<div class="overflow-x-auto rounded-3xl border border-base-300 bg-base-100 shadow-sm">
	<table class="table">
		<thead>
			<tr class="text-base-content/70">
				<th>ペア</th>
				<th>テンプレート</th>
				<th>状態</th>
				<th>実施予定日</th>
				<th class="text-right">操作</th>
			</tr>
		</thead>
		<tbody>
			{#each items as item (item.uuid)}
				<tr class="hover:bg-base-200/40">
					<td class="whitespace-nowrap">{pairLabel(item)}</td>
					<td>{item.template?.title || '1on1'}</td>
					<td>
						<span class="badge badge-sm {statusBadge(item.status ?? 0)}">
							{statusLabel(item.status ?? 0)}
						</span>
					</td>
					<td class="whitespace-nowrap">{formatDate(item.implementedAt)}</td>
					<td class="text-right">
						<div class="flex flex-wrap justify-end gap-1">
							{#if showAnswer}
								<a href="/meeting/{item.uuid}/answer" class="btn btn-xs btn-primary">記入</a>
							{/if}
							<a href="/meeting/{item.uuid}" class="btn btn-xs btn-ghost">閲覧</a>
							{#if showSetting}
								<a href="/meeting/{item.uuid}/setting" class="btn btn-xs btn-ghost">設定</a>
							{/if}
						</div>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
	{#if items.length === 0}
		<p class="p-8 text-center text-sm text-base-content/50">面談はありません</p>
	{/if}
</div>
