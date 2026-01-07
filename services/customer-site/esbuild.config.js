const esbuild = require('esbuild');
const fs = require('fs');
const path = require('path');

const entryDir = './web/dev/scripts/entrypoints';
const outDir = './web/public/js/entrypoints';

// собираем entry points автоматически
const entryPoints = Object.fromEntries(
	fs
		.readdirSync(entryDir)
		.filter(f => f.endsWith('.ts'))
		.map(f => [path.basename(f, '.ts'), path.join(entryDir, f)])
);

try {
	fs.rmSync(outDir, {
		recursive: true,
		force: true,
	});
} catch (error) {
	console.error(error);
	process.exit(1);
}

esbuild
	.build({
		entryPoints,
		outdir: outDir,

		bundle: true,
		format: 'esm',
		target: 'es2020',
		splitting: true,

		sourcemap: false,
		minify: true,
		treeShaking: true,

		entryNames: '[name]-[hash]',
		chunkNames: 'chunks/[name]-[hash]',
		metafile: true,

		loader: {
			'.ts': 'ts',
		},

		// define: {
		// 	__DEV__: JSON.stringify(isDev),
		// },
	})
	.then(result => {
		fs.writeFileSync(
			path.join(entryDir, 'manifest.json'),
			JSON.stringify(result.metafile.outputs, null, 2)
		);
	})
	.catch(() => process.exit(1));
