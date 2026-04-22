const form = document.getElementById("product-form");
const input = document.getElementById("product-id");
const submitButton = document.getElementById("submit-button");
const loadingState = document.getElementById("loading-state");
const errorState = document.getElementById("error-state");
const resultCard = document.getElementById("result-card");
const resultTitle = document.getElementById("result-title");
const cachePill = document.getElementById("cache-pill");
const productDetails = document.getElementById("product-details");
const marketDetails = document.getElementById("market-details");
const metaDetails = document.getElementById("meta-details");

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const productID = input.value.trim().toUpperCase();
  if (!productID) {
    renderError({
      error: "Please enter a product ID such as BTC-USD.",
    });
    return;
  }

  setLoading(true);
  clearError();

  try {
    const response = await fetch(`/products/${encodeURIComponent(productID)}`);
    const body = await response.json();

    if (!response.ok) {
      renderError(body);
      resultCard.classList.add("hidden");
      return;
    }

    renderProduct(body);
  } catch (error) {
    renderError({
      error: "Unable to reach the backend service.",
      details: error instanceof Error ? error.message : "",
    });
    resultCard.classList.add("hidden");
  } finally {
    setLoading(false);
  }
});

function setLoading(isLoading) {
  submitButton.disabled = isLoading;
  loadingState.classList.toggle("hidden", !isLoading);
  submitButton.textContent = isLoading ? "Loading..." : "Fetch Product";
}

function clearError() {
  errorState.textContent = "";
  errorState.classList.add("hidden");
}

function renderError(error) {
  const parts = [error.error || "Request failed."];
  if (error.details) {
    parts.push(error.details);
  }
  errorState.textContent = parts.join(" ");
  errorState.classList.remove("hidden");
}

function renderProduct(product) {
  clearError();
  resultCard.classList.remove("hidden");
  resultTitle.textContent = `${product.product_name} (${product.product_id})`;
  cachePill.textContent = `cache: ${product.cache_status}`;

  renderDetails(productDetails, [
    ["Product ID", product.product_id],
    ["Market Pair", product.market_pair],
    ["Product Name", product.product_name],
    ["Base Currency", product.base_currency],
    ["Quote Currency", product.quote_currency],
    ["Status", product.status],
    ["Trading Enabled", String(product.is_trading_enabled)],
  ]);

  renderDetails(marketDetails, [
    ["Price", product.price],
    ["24H Price Change", product.price_change_24h],
  ]);

  renderDetails(metaDetails, [
    ["Cache Status", product.cache_status],
    ["Retrieved At", product.retrieved_at],
    ["Source", product.source],
  ]);
}

function renderDetails(container, items) {
  container.innerHTML = "";

  for (const [label, value] of items) {
    const wrapper = document.createElement("div");
    const term = document.createElement("dt");
    const description = document.createElement("dd");

    term.textContent = label;
    description.textContent = value || "-";

    wrapper.appendChild(term);
    wrapper.appendChild(description);
    container.appendChild(wrapper);
  }
}
