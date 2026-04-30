document.addEventListener("DOMContentLoaded", function () {
  const ASSET = "/EpicFreeGamesList/assets/";

  // Platform options definition
  const PLATFORMS = [
    {
      value: "",
      label: "All",
      icons: [
        { src: ASSET + "windows.svg", alt: "PC" },
        { src: ASSET + "android.svg", alt: "Android" },
        { src: ASSET + "apple.svg",   alt: "iOS" },
      ],
    },
    {
      value: "pc",
      label: "PC",
      icons: [{ src: ASSET + "windows.svg", alt: "PC" }],
    },
    {
      value: "android",
      label: "Android",
      icons: [{ src: ASSET + "android.svg", alt: "Android" }],
    },
    {
      value: "ios",
      label: "iOS",
      icons: [{ src: ASSET + "apple.svg", alt: "iOS" }],
    },
  ];

  // Currently selected platform value ("" = All)
  let selectedPlatform = "";

  // Build icons HTML for a platform option
  function iconsHtml(icons) {
    return icons
      .map(i => `<img src="${i.src}" alt="${i.alt}">`)
      .join("");
  }

  // Build the full selected-value inner HTML
  function selectedHtml(option) {
    return `${iconsHtml(option.icons)}<span class="pf-label">${option.label}</span><span class="pf-arrow">&#9660;</span>`;
  }

  // Register custom DataTables row filter (runs on every draw)
  $.fn.dataTable.ext.search.push(function (settings, searchData, dataIndex, rowData) {
    if (selectedPlatform === "") return true;
    const rowPlatform = (rowData.platform || "pc").toLowerCase();
    return rowPlatform === selectedPlatform;
  });

  fetch("/EpicFreeGamesList/epic_free_games.json")
    .then(response => response.json())
    .then(data => {
      const dataTable = new DataTable(document.getElementById("epicGamesTable"), {
        data: data,
        pageLength: 100,
        order: [
          [2, "desc"],
          [3, "desc"],
        ],
        columns: [
          {
            title: "Title",
            data: "gameTitle",
            render: function (data, type, row) {
              if (type === "display") {
                return `<a href="${row.epicStoreLink}" target="_blank">${data}</a>`;
              }
              return data;
            },
          },
          {
            title: "Platform",
            data: "platform",
            render: function (data, type) {
              const platform = data ? data.toLowerCase() : "pc";
              let iconPath = "";
              let label = "";
              switch (platform) {
                case "android":
                  iconPath = ASSET + "android.svg";
                  label = "Android";
                  break;
                case "ios":
                  iconPath = ASSET + "apple.svg";
                  label = "iOS";
                  break;
                case "pc":
                default:
                  iconPath = ASSET + "windows.svg";
                  label = "PC";
                  break;
              }
              return `<img src="${iconPath}" alt="${label}" title="${label}" style="width:20px; height:20px;">`;
            },
          },
          {
            title: "Date",
            data: "freeDate",
          },
          {
            title: "<img src='/EpicFreeGamesList/assets/steamdb.svg' alt='SteamDB' title='SteamDB Rating' style='width:14px;height:14px;vertical-align:middle;margin-right:4px;'>SteamDB",
            data: "steamDBRating",
            render: function (data, type, row) {
              if (type === "display") {
                if (!data) return "";
                const label = data + "%";
                if (row.steamDBUrl) {
                  return `<a href="${row.steamDBUrl}" target="_blank">${label}</a>`;
                }
                return label;
              }
              return data || 0;
            },
          },
          {
            title: "<img src='/EpicFreeGamesList/assets/metacritic.ico' alt='Metacritic' title='Metascore' style='width:14px;height:14px;vertical-align:middle;margin-right:4px;'>Metascore",
            data: "metacriticScore",
            render: function (data, type, row) {
              if (type === "display") {
                if (!data) return "";
                if (row.metacriticUrl) {
                  return `<a href="${row.metacriticUrl}" target="_blank">${data}</a>`;
                }
                return data;
              }
              return data || 0;
            },
          },
          {
            title: "<img src='/EpicFreeGamesList/assets/epicgames.ico' alt='Epic Games' title='Epic Rating' style='width:14px;height:14px;vertical-align:middle;margin-right:4px;'>Epic Rating",
            data: "epicRating",
          },
        ],

        initComplete: function () {
          // Build dropdown HTML
          const menuItems = PLATFORMS.map(opt => {
            const val = opt.value === "" ? "__all__" : opt.value;
            return `<li data-platform="${opt.value}">${iconsHtml(opt.icons)}<span>${opt.label}</span></li>`;
          }).join("");

          const $dropdown = $(`
            <div class="platform-filter" id="platform-filter">
              <div class="platform-filter-selected" id="pf-selected">
                ${selectedHtml(PLATFORMS[0])}
              </div>
              <ul class="platform-filter-menu" id="pf-menu">
                ${menuItems}
              </ul>
            </div>
          `);

          // Wrap the DataTables filter label+input in a flex toolbar
          const $dtFilter = $("#epicGamesTable_filter");
          const $toolbar = $('<div id="dt-toolbar"></div>');
          $dtFilter.before($toolbar);
          $toolbar.append($dropdown).append($dtFilter);

          // Toggle menu open/closed
          $("#pf-selected").on("click", function (e) {
            e.stopPropagation();
            $("#pf-menu").toggleClass("open");
          });

          // Select an option
          $("#pf-menu").on("click", "li", function () {
            const val = $(this).data("platform");
            selectedPlatform = val;
            const match = PLATFORMS.find(p => p.value === val) || PLATFORMS[0];
            $("#pf-selected").html(selectedHtml(match));
            $("#pf-menu").removeClass("open");
            dataTable.draw();
          });

          // Close on outside click
          $(document).on("click", function () {
            $("#pf-menu").removeClass("open");
          });
        },
      });
    })
    .catch(error => {
      console.error("Error loading data:", error);
    });
});