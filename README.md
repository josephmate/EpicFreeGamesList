
```shell
go get -u github.com/chromedp/chromedp
```
chromedp was unsuccessful since cloudflare gives us a challenge which we cannot answer.

instead I found this stack overflow
https://stackoverflow.com/questions/75594997/is-it-possible-to-get-data-about-specific-game-from-epic-games-store

I was able to run the query without issue in powershell:
```powershell
$response = Invoke-WebRequest -Uri 'https://graphql.epicgames.com/graphql?query=query%20searchStoreQuery(%20$allowCountries:%20String%20$category:%20String%20$namespace:%20String%20$itemNs:%20String%20$sortBy:%20String%20$sortDir:%20String%20$start:%20Int%20$tag:%20String%20$releaseDate:%20String%20$withPrice:%20Boolean%20=%20true%20)%20{%20Catalog%20{%20searchStore(%20allowCountries:%20$allowCountries%20category:%20$category%20count:%201%20country:%20%22US%22%20keywords:%20%22celeste%22%20namespace:%20$namespace%20itemNs:%20$itemNs%20sortBy:%20$sortBy%20sortDir:%20$sortDir%20releaseDate:%20$releaseDate%20start:%20$start%20tag:%20$tag%20)%20{%20elements%20{%20title%20description%20keyImages%20{%20type%20url%20}%20seller%20{%20name%20}%20categories%20{%20path%20}%20price(country:%20%22US%22)%20@include(if:%20$withPrice)%20{%20totalPrice%20{%20fmtPrice(locale:%20%22en-US%22)%20{%20discountPrice%20}%20}%20}%20}%20}%20}%20}'
Write-Output $response.StatusCode
200
Write-Output $response.Content
{"data":{"Catalog":{"searchStore":{"elements":[{"title":"Celeste","description":"Help Madeline survive her inner demons on her journey to the top of Celeste Mountain, in this super-tight platformer from the creators of TowerFall. Brave hundreds of hand-crafted challenges, uncover devious secrets, and piece together the mystery of the mountain.","keyImages":[{"type":"OfferImageWide","url":"https://cdn1.epicgames.com/b671fbc7be424e888c9346a9a6d3d9db/offer/Celeste - landscape offer image-2560x1440-0b9b94fd493d817704ecfdf4c704989a.jpg"},{"type":"OfferImageTall","url":"https://cdn1.epicgames.com/offer/b671fbc7be424e888c9346a9a6d3d9db/CodeRedemption_Celeste-340x440-873dacb76be8c59d659757b4b5284a611_1200x1600-9e39079f1ce19738e65d349a8bc98650"},{"type":"Thumbnail","url":"https://cdn1.epicgames.com/offer/b671fbc7be424e888c9346a9a6d3d9db/CodeRedemption_Celeste-340x440-873dacb76be8c59d659757b4b5284a611_1200x1600-9e39079f1ce19738e65d349a8bc98650"}],"seller":{"name":"Maddy Makes Games"},"categories":[{"path":"games"},{"path":"games/edition/base"},{"path":"games/edition"},{"path":"applications"}],"price":{"totalPrice":{"fmtPrice":{"discountPrice":"$19.99"}}}}]}}},"extensions":{}}
```

next step is to explore their graphql API and look for the rating some where.