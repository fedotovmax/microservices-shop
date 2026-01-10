const SVGSprite = require('svg-sprite');
const fs = require('fs');
const path = require('path');
const crypto = require('crypto');

const ICONS_PATH = './web/dev/icons';
const OUTPUT_PATH = './web/public/icons';

const GO_ICON_IDS_PATH = './internal/components/icon/ids.go';
const GO_SPRITE_PATH = './internal/resources/sprite.go';

const config = {
	mode: {
		symbol: {
			sprite: 'sprite.svg',
		},
	},
	shape: {
		id: {
			generator: 'icon-%s',
		},
	},
};

const sprite = new SVGSprite(config);
const icons = [];

if (fs.existsSync(OUTPUT_PATH)) {
	fs.readdirSync(OUTPUT_PATH).forEach(file => {
		fs.unlinkSync(path.join(OUTPUT_PATH, file));
	});
} else {
	fs.mkdirSync(OUTPUT_PATH, { recursive: true });
}

fs.readdirSync(ICONS_PATH).forEach(file => {
	if (!file.endsWith('.svg')) return;

	const name = path.basename(file, '.svg');
	icons.push(name);

	const fullPath = path.join(ICONS_PATH, file);
	sprite.add(fullPath, file, fs.readFileSync(fullPath, 'utf8'));
});

sprite.compile((err, result) => {
	if (err) throw err;

	let svg = result.symbol.sprite.contents.toString();

	svg = svg.replace(/\sfill="[^"]*"/g, '');

	const hash = crypto.createHash('md5').update(svg).digest('hex').slice(0, 8);

	const spriteName = `sprite.${hash}.svg`;

	// SVG sprite
	fs.mkdirSync(OUTPUT_PATH, { recursive: true });
	fs.writeFileSync(path.join(OUTPUT_PATH, spriteName), svg);

	// Go: icon IDs
	const goIconIDs = `package icon

type ID string

const (
${icons
	.sort()
	.map(name => `\tIcon${toGoIdent(name)} ID = "icon-${name}"`)
	.join('\n')}
)
`;

	fs.mkdirSync(path.dirname(GO_ICON_IDS_PATH), { recursive: true });
	fs.writeFileSync(GO_ICON_IDS_PATH, goIconIDs);

	// Go: sprite name
	const goSprite = `package resources

const SpriteName = "/public/icons/${spriteName}"
`;

	fs.mkdirSync(path.dirname(GO_SPRITE_PATH), { recursive: true });
	fs.writeFileSync(GO_SPRITE_PATH, goSprite);

	console.log(`✔ SVG sprite generated: ${spriteName}`);
	console.log(`✔ Go icon IDs generated`);
	console.log(`✔ Go sprite name generated`);
});

function toGoIdent(name) {
	return name
		.replace(/[-_]+(.)?/g, (_, c) => (c ? c.toUpperCase() : ''))
		.replace(/^(.)/, c => c.toUpperCase());
}
