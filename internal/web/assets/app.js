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

let analyzeRequestId = 0;

function renderDensity(rows) {
  densityBodyEl.replaceChildren();
  densitySummaryEl.textContent = `${rows.length} unique words`;

  if (rows.length === 0) {
    densityEmptyEl.hidden = false;
    return;
  }

  densityEmptyEl.hidden = true;

  const fragment = document.createDocumentFragment();

  for (const rowData of rows) {
    const row = document.createElement("tr");

    const wordCell = document.createElement("td");
    wordCell.textContent = rowData.word;

    const countCell = document.createElement("td");
    countCell.textContent = String(rowData.count);

    const densityCell = document.createElement("td");
    densityCell.textContent = rowData.density;

    row.append(wordCell, countCell, densityCell);
    fragment.appendChild(row);
  }

  densityBodyEl.appendChild(fragment);
}

async function requestAnalysis(text) {
  const response = await fetch("/api/analyze", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ text }),
  });

  if (!response.ok) {
    throw new Error(`analysis request failed with status ${response.status}`);
  }

  return response.json();
}

async function updateUI() {
  const requestId = ++analyzeRequestId;
  const analysis = await requestAnalysis(editor.value);

  if (requestId !== analyzeRequestId) {
    return;
  }

  charactersEl.textContent = String(analysis.characters);
  wordsEl.textContent = String(analysis.words);
  sentencesEl.textContent = String(analysis.sentences);
  paragraphsEl.textContent = String(analysis.paragraphs);
  spacesEl.textContent = String(analysis.spaces);
  renderDensity(analysis.densityRows);
}

function triggerUIUpdate() {
  void updateUI().catch(() => {});
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
  triggerUIUpdate();
});

themeToggle.addEventListener("click", () => {
  const nextTheme = document.body.dataset.theme === "dark" ? "light" : "dark";
  setTheme(nextTheme);
  saveThemeState(nextTheme);
});

loadTextState();
setTheme(loadTheme());
triggerUIUpdate();
