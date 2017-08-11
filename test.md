<script src="https://ajax.googleapis.com/ajax/libs/jquery/1.11.1/jquery.min.js">
</script>

<p id="p1">https://github.com/jvanderl/flogo-components/activity/blewrite</p>

<button onclick="copyToClipboard('#p1')">Copy to clipboard</button>

<script>
function copyToClipboard(element) {
  var $temp = $("<input>");
  $("body").append($temp);
  $temp.val($(element).text()).select();
  document.execCommand("copy");
  $temp.remove();
}

//# sourceURL=pen.js
</script>
