+++
title = "<% title %>"
date = "<% date %>"
coverImage = "<% thumbnail %>"
thumbnailImage = "<% thumbnail %>"
coverMeta = "in"
thumbnailImagePosition = "right"
categories = ["<% category %>"]
tags = [<% tags %>]
comments = false
+++

<div class="brunch">
<% content %>
</div>

<div id="notice" style="display: none;">
{{< alert info >}}
이 글은 <a href="<% link %>">브런치 글</a>에서 자동으로 생성되었습니다.
{{< /alert >}}
</div>

<script id="show-notification">
;(() => {
  window.addEventListener('load', () => {
    ~document.querySelector('meta[name=author]').content.indexOf('Wooseop') &&
      (document.getElementById('notice').style.display = 'initial') &&
      document.getElementById('show-notification').remove()
  })
})();
</script>
