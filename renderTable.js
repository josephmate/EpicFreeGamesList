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
                [1, 'desc'],
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