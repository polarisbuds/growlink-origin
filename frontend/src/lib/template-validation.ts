export type TemplateTopicDraft = { question?: string };

export function validateTemplateForm(
	title: string,
	description: string,
	topics: TemplateTopicDraft[],
): string | null {
	if (!title.trim()) return 'タイトルを入力してください';
	if (!description.trim()) return '説明を入力してください';
	if (topics.length === 0) return '質問を1件以上追加してください';
	for (let i = 0; i < topics.length; i++) {
		if (!topics[i].question?.trim()) {
			return `Q${i + 1} の質問文を入力してください`;
		}
	}
	return null;
}

export function formatTemplateCreatedAt(iso?: string): string {
	if (!iso) return '—';
	const d = new Date(iso);
	if (Number.isNaN(d.getTime())) return '—';
	return d.toLocaleDateString('ja-JP');
}
