document.addEventListener("DOMContentLoaded", function () {
  fetch("/EpicFreeGamesList/epic_free_games.json") // Replace with the actual URL of your JSON file
    .then(response => response.json())
    .then(data => {
      // Initialize DataTable with the retrieved data
      const dataTable = new DataTable(document.getElementById("epicGamesTable"), {
        data: data,
        pageLength: 100,
        order: [
          [2, 'desc'],
          [3, 'desc'],
        ],
        columns: [
          {
            title: "Title",
            data: "gameTitle",
            render: function (data, type, row, meta) {
              if (type === "display") {
                return `<a href="${row.epicStoreLink}" target="_blank">${data}</a>`;
              }
              return data;
            }
          },
          {
            title: "Platform",
            data: "platform",
            render: function (data, type, row, meta) {
              // Default to PC if missing
              let platform = data ? data.toLowerCase() : "pc";
              let iconPath = "";
              let label = "";

              switch (platform) {
                case "android":
                  iconPath = "/EpicFreeGamesList/assets/android.svg";
                  label = "Android";
                  break;
                case "ios":
                  iconPath = "/EpicFreeGamesList/assets/apple.svg";
                  label = "iOS";
                  break;
                case "pc":
                default:
                  iconPath = "/EpicFreeGamesList/assets/windows.svg";
                  label = "PC";
                  break;
              }

              return `<img src="${iconPath}" alt="${label}" title="${label}" style="width:20px; height:20px;">`;
            }
          },
          {
            title: "Date",
            data: "freeDate"
          },
          {
            title: "Epic Rating",
            data: "epicRating"
          },
        ]
      });
    })
    .catch(error => {
      console.error("Error loading data:", error);
    });
});