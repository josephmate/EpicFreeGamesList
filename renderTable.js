document.addEventListener("DOMContentLoaded", function () {
  fetch("/epic_free_games.json") // Replace with the actual URL of your JSON file
      .then(response => response.json())
      .then(data => {
          // Initialize DataTable with the retrieved data
          const dataTable = new DataTable(document.getElementById("epicGamesTable"), {
              data: data,
              columns: [
                  {
                      data: "gameTitle",
                      render: function (data, type, row, meta) {
                          if (type === "display") {
                              return `<a href="${row.epicStoreLink}" target="_blank">${data}</a>`;
                          }
                          return data;
                      }
                  },
                  { data: "freeDate" }
              ]
          });
      })
      .catch(error => {
          console.error("Error loading data:", error);
      });
});