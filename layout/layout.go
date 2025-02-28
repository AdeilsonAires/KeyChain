package layout

import (
	database "KeyChain/dataBase"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type LayoutWindow struct {
	App                     fyne.App
	Win                     fyne.Window
	MainContainer           *fyne.Container
	Toolbar                 *fyne.Container
	ContentContainer        *fyne.Container
	ManagePasswordsButton   *widget.Button
	CreateNewPasswordButton *widget.Button
	OpenWindows []fyne.Window
}

func (hw *LayoutWindow) showEditPasswordDialog(cred database.Credential) {
	editWin := hw.App.NewWindow("Editar Senha")

	senhaEntry := widget.NewEntry()
	senhaEntry.SetText(cred.Senha)
	senhaEntry.OnChanged = func(text string) {
		if len(text) > 15 {
			senhaEntry.SetText(text[:15])
			senhaEntry.CursorColumn = 15
		}
	}

	btnGeneratePassword := widget.NewButton("Gerar Nova Senha", func() {
		senhaEntry.SetText(gerarSenha(12))
	})

	btnSave := widget.NewButton("Salvar Alterações", func() {
		newPassword := senhaEntry.Text

		if newPassword == "" {
			dialog.ShowInformation("Erro", "A senha não pode estar vazia!", editWin)
			return
		}

		err := database.UpdatePassword(cred.ID, newPassword)
		if err != nil {
			dialog.ShowError(err, editWin)
			return
		}

		dialog.ShowInformation("Sucesso", "Senha atualizada com sucesso!", editWin)
		editWin.Close()
		hw.SetMainContainer("managePassword")
		hw.Win.SetContent(hw.MainContainer)
	})

	form := container.NewVBox(
		widget.NewLabel("Nova Senha:"),
		senhaEntry,
		btnGeneratePassword,
		btnSave,
	)

	editWin.SetContent(form)
	editWin.Resize(fyne.NewSize(400, 200))
	editWin.Show()
}

// Função para gerar senha aleatória
func gerarSenha(tamanho int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"
	senha := make([]byte, tamanho)

	for i := range senha {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		senha[i] = charset[num.Int64()]
	}
	return string(senha)
}

func showExitMessage(win fyne.Window) {
	dialog.ShowInformation("Atenção", "O aplicativo será fechado em 3 segundos...", win)
	time.AfterFunc(3*time.Second, func() {
		win.Close()
	})
}

func (hw *LayoutWindow) SetMainContainer(win string) {
	hw.ContentContainer = hw.setContentContainer(win)
	hw.Toolbar = hw.setToolbar(win)
	if win == "home" {
		hw.MainContainer = container.NewBorder(hw.Toolbar, hw.ContentContainer, nil, nil, nil)
	} else if win == "managePassword" {
		hw.MainContainer = container.NewBorder(hw.Toolbar, nil, nil, nil, hw.ContentContainer)
	}
}

func loadImageFromAsset() (*fyne.StaticResource, error) {
	return fyne.NewStaticResource("keychain_logo.png", resourceKeychainlogoPng.Content()), nil
}

func (hw *LayoutWindow) setToolbar(win string) *fyne.Container {

	var toolbarLeft *widget.Toolbar
	if win == "home" {
		toolbarLeft = widget.NewToolbar(
			widget.NewToolbarAction(theme.InfoIcon(), func() {
				infoDialog := dialog.NewInformation("Sobre", "Aplicativo de Gerenciamento de Senhas", hw.Win)
				infoDialog.Show()
				time.AfterFunc(3*time.Second, func() {
					infoDialog.Hide()
				})
			}),
		)
	} else if win == "managePassword" {
		toolbarLeft = widget.NewToolbar(
			widget.NewToolbarAction(theme.NavigateBackIcon(), func() {
				hw.SetMainContainer("home")
				hw.Win.SetContent(hw.MainContainer)
			}),
		)
	}

	toolbarRight := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentClearIcon(), func() {
			showExitMessage(hw.Win)
		}),
	)

	logoResource, err := loadImageFromAsset()
	var logo *canvas.Image
	if err != nil {
		fmt.Println("Erro ao carregar imagem:", err)
		logo = canvas.NewImageFromResource(theme.FyneLogo())
	} else {
		logo = canvas.NewImageFromResource(logoResource)
		logo.FillMode = canvas.ImageFillContain
		logo.SetMinSize(fyne.NewSize(200, 100))
	}

	hw.Toolbar = container.NewHBox(
		toolbarLeft,
		layout.NewSpacer(),
		container.NewCenter(logo),
		layout.NewSpacer(),
		toolbarRight,
	)

	return hw.Toolbar
}

func (hw *LayoutWindow) setContentContainer(win string) *fyne.Container {
	credentialTable := &widget.Table{}
	if win == "home" {
		hw.ManagePasswordsButton = widget.NewButtonWithIcon("Gerenciar Suas Senhas", theme.FolderOpenIcon(), func() {
			hw.SetMainContainer("managePassword")
			hw.Win.SetContent(hw.MainContainer)
		})

		hw.CreateNewPasswordButton = widget.NewButtonWithIcon("Cadastrar Nova Senha", theme.DocumentCreateIcon(), func() {
			hw.setFormNewPassword(hw.App)
		})

		hw.ManagePasswordsButton.Alignment = widget.ButtonAlignCenter
		hw.ManagePasswordsButton.Importance = widget.SuccessImportance

		hw.CreateNewPasswordButton.Alignment = widget.ButtonAlignCenter
		hw.CreateNewPasswordButton.Importance = widget.SuccessImportance

		hw.ContentContainer = container.NewVBox(
			hw.ManagePasswordsButton,
			hw.CreateNewPasswordButton,
		)

		hw.Win.Resize(fyne.NewSize(250, 250))

	} else if win == "managePassword" {
		list, err := database.GetCredentials()
		if err != nil {
			dialog.ShowError(err, hw.Win)
			return nil
		}
		var senhaVisivel = make(map[int]bool)
		credentialTable = widget.NewTable(
			func() (int, int) { return len(list) + 1, 7 }, 
			func() fyne.CanvasObject {
				label := widget.NewLabel("")
				button := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), nil)
				button.Hide()
				return container.NewMax(label, button)
			},
			func(tableCell widget.TableCellID, cell fyne.CanvasObject) {
				cont := cell.(*fyne.Container)
				label := cont.Objects[0].(*widget.Label)
				button := cont.Objects[1].(*widget.Button)
				label.Show()
				button.Hide()
		
				if tableCell.Row == 0 {
					switch tableCell.Col {
					case 0:
						label.SetText("Site/App")
					case 1:
						label.SetText("Usuário")
					case 2:
						label.SetText("Senha")
					case 3:
						label.SetText("Mostrar")
					case 4:
						label.SetText("Copiar")
					case 5:
						label.SetText("Editar")
					case 6:
						label.SetText("Deletar")
					}
				} else {
					index := tableCell.Row - 1
					switch tableCell.Col {
					case 0:
						label.SetText(list[index].SiteApp)
					case 1:
						label.SetText(list[index].Usuario)
					case 2:
						if senhaVisivel[index] {
							label.SetText(list[index].Senha)
						} else {
							label.SetText("**************")
						}
					case 3:
						label.Hide()
						buttonIndex := index // Criamos uma cópia local do índice para evitar capturas incorretas
						if senhaVisivel[buttonIndex] {
							button.SetIcon(theme.VisibilityOffIcon()) // Ícone de "esconder"
						} else {
							button.SetIcon(theme.VisibilityIcon()) // Ícone de "mostrar"
						}
						button.OnTapped = func() {
							senhaVisivel[buttonIndex] = !senhaVisivel[buttonIndex] // Alterna o estado
							credentialTable.Refresh() // Atualiza a tabela para refletir a mudança
						}
						button.Show()
					case 4:
						label.Hide()
						button.SetIcon(theme.ContentCopyIcon())
						button.OnTapped = func() {
							clipboard := hw.Win.Clipboard()
							clipboard.SetContent(list[index].Senha)
							dialog.ShowInformation("Copiado", "Senha copiada para a área de transferência!", hw.Win)
						}
						button.Show()
					case 5:
						label.Hide()
						button.SetIcon(theme.DocumentCreateIcon())
						button.OnTapped = func() {
							hw.showEditPasswordDialog(list[index])
						}
						button.Show()
					case 6:
						label.Hide()
						button.SetIcon(theme.DeleteIcon())
						button.OnTapped = func() {
							dialog.NewConfirm(
								"Confirmação",
								"Deseja realmente excluir essa senha?",
								func(confirm bool) {
									if confirm {
										err := database.DeleteCredential(list[index].ID)
										if err != nil {
											dialog.ShowError(err, hw.Win)
											return
										}
										hw.SetMainContainer("managePassword")
										hw.Win.SetContent(hw.MainContainer)
										dialog.ShowInformation("Sucesso", "Senha excluída com sucesso!", hw.Win)
									}
								},
								hw.Win,
							).Show()
						}
						button.Show()
					}
				}
			},
		)
		
		for row := 0; row < len(list)+1; row++ {
			credentialTable.SetRowHeight(row, 40)
		}
		credentialTable.SetColumnWidth(0, 170)
		credentialTable.SetColumnWidth(1, 250)
		credentialTable.SetColumnWidth(2, 150)
		credentialTable.SetColumnWidth(3, 60)
		credentialTable.SetColumnWidth(4, 60)
		credentialTable.SetColumnWidth(5, 60)
		credentialTable.SetColumnWidth(6, 60)

		hw.ContentContainer = container.NewBorder(
			nil,
			nil,
			nil,
			nil,
			container.NewMax(credentialTable),
		)
		hw.Win.Resize(fyne.NewSize(845, 450))
		return hw.ContentContainer
	}
	return hw.ContentContainer
}

func (hw *LayoutWindow) setFormNewPassword(app fyne.App) {
	formWin := app.NewWindow("Cadastro de Senha")

	siteAppEntry := widget.NewEntry()
	siteAppEntry.SetPlaceHolder("Digite o nome do site/app (Tamanho maximo 20 caracteres)")
	siteAppEntry.OnChanged = func(text string) {
	if len(text) > 20 {
		siteAppEntry.SetText(text[:20])
		siteAppEntry.CursorColumn = 20 
	}
}
	usuarioEntry := widget.NewEntry()
	usuarioEntry.SetPlaceHolder("Digite seu usuário (Tamanho maximo 30 caracteres)")
	usuarioEntry.OnChanged = func(text string) {
		if len(text) > 30 {
			usuarioEntry.SetText(text[:30])
			usuarioEntry.CursorColumn = 30 
		}
	}
	senhaEntry := widget.NewEntry()
	senhaEntry.SetPlaceHolder("Digite sua senha (Tamanho maximo 15 caracteres)")
	senhaEntry.OnChanged = func(text string) {
		if len(text) > 15 {
			senhaEntry.SetText(text[:15])
			senhaEntry.CursorColumn = 15
		}
	}

	btnGeneratePassword := widget.NewButton("Gerar Senha", func() {
		senhaEntry.SetText(gerarSenha(15))
	})

	btnSave := widget.NewButton("Salvar", func() {
		websiteApp := siteAppEntry.Text
		user := usuarioEntry.Text
		password := senhaEntry.Text

		if websiteApp == "" || user == "" || password == "" {
			dialog.ShowInformation("Erro", "Todos os campos devem ser preenchidos!", formWin)
			return
		}
		err := database.InsertCredential(websiteApp, user, password)
		if err != nil {
			dialog.ShowError(err, formWin)
			return
		}

		dialog.ShowInformation("Sucesso", fmt.Sprintf("Site: %s\nUsuário: %s\nSenha salva!", websiteApp, user), formWin)

		siteAppEntry.SetText("")
		usuarioEntry.SetText("")
		senhaEntry.SetText("")
	})

	form := container.NewVBox(
		widget.NewLabel("Site/App:"),
		siteAppEntry,
		widget.NewLabel("Usuário:"),
		usuarioEntry,
		widget.NewLabel("Senha:"),
		senhaEntry,
		btnGeneratePassword,
		btnSave,
	)

	formWin.SetContent(form)
	formWin.Resize(fyne.NewSize(450, 300))
	formWin.Show()
}
