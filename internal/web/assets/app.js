const TEXT_STORAGE_KEY = "charcount_text";
const THEME_STORAGE_KEY = "charcount_theme";

const editor = document.getElementById("editor");
const themeToggle = document.getElementById("theme-toggle");
const charactersEl = document.getElementById("characters");
const wordsEl = document.getElementById("words");
const sentencesEl = document.getElementById("sentences");
const paragraphsEl = document.getElementById("paragraphs");
const spacesEl = document.getElementById("spaces");
const densitySummaryEl = document.getElementById("density-summary");
const densityBodyEl = document.getElementById("density-body");
const densityEmptyEl = document.getElementById("density-empty");
const defaultText = editor.value;

function extractWords(text) {
  const words = text.toLowerCase().match(/\p{L}[\p{L}\p{N}'-]*/gu);
  return words || [];
}

function countSentences(text) {
  const trimmed = text.trim();
  if (!trimmed) {
    return 0;
  }

  const matches = trimmed.match(/[^.!?]+[.!?]+|[^.!?]+$/g);
  if (!matches) {
    return 0;
  }

  return matches.map((sentence) => sentence.trim()).filter(Boolean).length;
}

function countParagraphs(text) {
  const normalized = text.replace(/\r\n?/g, "\n").trim();
  if (!normalized) {
    return 0;
  }

  return normalized.split(/\n\s*\n/).filter((paragraph) => paragraph.trim() !== "").length;
}

function buildDensityRows(words) {
  const counts = new Map();

  for (const word of words) {
    counts.set(word, (counts.get(word) || 0) + 1);
  }

  return Array.from(counts.entries()).sort((a, b) => {
    const byCount = b[1] - a[1];
    if (byCount !== 0) {
      return byCount;
    }

    return a[0].localeCompare(b[0]);
  });
}

function renderDensity(rows, wordCount) {
  densityBodyEl.replaceChildren();
  densitySummaryEl.textContent = `${rows.length} unique words`;

  if (rows.length === 0 || wordCount === 0) {
    densityEmptyEl.hidden = false;
    return;
  }

  densityEmptyEl.hidden = true;

  const fragment = document.createDocumentFragment();

  for (const [word, count] of rows) {
    const density = ((count / wordCount) * 100).toFixed(2);
    const row = document.createElement("tr");

    const wordCell = document.createElement("td");
    wordCell.textContent = word;

    const countCell = document.createElement("td");
    countCell.textContent = String(count);

    const densityCell = document.createElement("td");
    densityCell.textContent = `${density}%`;

    row.append(wordCell, countCell, densityCell);
    fragment.appendChild(row);
  }

  densityBodyEl.appendChild(fragment);
}

function analyzeText(text) {
  const words = extractWords(text);
  const wordCount = words.length;
  const result = {
    characters: text.length,
    words: wordCount,
    sentences: countSentences(text),
    paragraphs: countParagraphs(text),
    spaces: (text.match(/ /g) || []).length,
    densityRows: buildDensityRows(words),
  };

  return result;
}

function updateUI() {
  const analysis = analyzeText(editor.value);
  charactersEl.textContent = String(analysis.characters);
  wordsEl.textContent = String(analysis.words);
  sentencesEl.textContent = String(analysis.sentences);
  paragraphsEl.textContent = String(analysis.paragraphs);
  spacesEl.textContent = String(analysis.spaces);
  renderDensity(analysis.densityRows, analysis.words);
}

function saveTextState() {
  try {
    localStorage.setItem(TEXT_STORAGE_KEY, editor.value);
  } catch (_error) {
    // Ignore storage failures and keep editing available.
  }
}

function loadTextState() {
  try {
    const savedText = localStorage.getItem(TEXT_STORAGE_KEY);
    if (savedText !== null) {
      editor.value = savedText;
      return;
    }
  } catch (_error) {
    // Ignore storage failures and use fallback text.
  }

  editor.value = defaultText;
}

function saveThemeState(theme) {
  try {
    localStorage.setItem(THEME_STORAGE_KEY, theme);
  } catch (_error) {
    // Ignore storage failures and keep rendering available.
  }
}

function setTheme(theme) {
  document.body.dataset.theme = theme;
  themeToggle.setAttribute("aria-pressed", String(theme === "dark"));
  themeToggle.textContent = theme === "dark" ? "Light mode" : "Dark mode";
}

function loadTheme() {
  try {
    const storedTheme = localStorage.getItem(THEME_STORAGE_KEY);
    if (storedTheme === "dark" || storedTheme === "light") {
      return storedTheme;
    }
  } catch (_error) {
    // Ignore storage failures and use fallback theme.
  }

  if (window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches) {
    return "dark";
  }

  return "light";
}

editor.addEventListener("input", () => {
  saveTextState();
  updateUI();
});

themeToggle.addEventListener("click", () => {
  const nextTheme = document.body.dataset.theme === "dark" ? "light" : "dark";
  setTheme(nextTheme);
  saveThemeState(nextTheme);
});

loadTextState();
setTheme(loadTheme());
updateUI();
