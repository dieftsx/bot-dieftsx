-- OBS Live Bot - Real-Time Typing Sync for Neovim
-- Adicione ao seu init.lua do Neovim para captura de digitação em tempo real (0ms de delay):
-- dofile(os.getenv("HOME") .. "/projects/bot-dieftsx/scripts/obs_nvim.lua")

local function sync_typing()
  local file = vim.api.nvim_buf_get_name(0)
  if file == "" or file:match("^term://") then
    return
  end

  local lines = vim.api.nvim_buf_get_lines(0, 0, 80, false)
  local content = table.concat(lines, "\n")
  local name = vim.fn.fnamemodify(file, ":t")

  local f = io.open("/tmp/obs_typing.txt", "w")
  if f then
    f:write(name .. "\n---OBS_SPLIT---\n" .. content)
    f:close()
  end
end

vim.api.nvim_create_autocmd({ "TextChangedI", "TextChanged", "BufEnter", "BufWritePost" }, {
  callback = sync_typing,
})
